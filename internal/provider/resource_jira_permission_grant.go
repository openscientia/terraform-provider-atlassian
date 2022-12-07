package atlassian

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/planmodifiers/stringmodifiers"
)

type (
	jiraPermissionGrantResource struct {
		p atlassianProvider
	}

	jiraPermissionGrantResourceModel struct {
		ID                 types.String                    `tfsdk:"id"`
		PermissionSchemeID types.String                    `tfsdk:"permission_scheme_id"`
		Holder             *jiraPermissionGrantHolderModel `tfsdk:"holder"`
		Permission         types.String                    `tfsdk:"permission"`
	}

	jiraPermissionGrantHolderModel struct {
		Type      types.String `tfsdk:"type"`
		Parameter types.String `tfsdk:"parameter"`
	}
)

var (
	_            resource.Resource                = (*jiraPermissionGrantResource)(nil)
	_            resource.ResourceWithImportState = (*jiraPermissionGrantResource)(nil)
	holder_types []string                         = []string{
		"anyone", "applicationRole", "assignee", "group", "groupCustomField", "projectLead",
		"projectRole", "reporter", "sd.customer.portal.only", "user", "userCustomField",
	}
	built_in_permissions []string = []string{
		// Project permissions
		"ADMINISTER_PROJECTS",
		"BROWSE_PROJECTS",
		"MANAGE_SPRINTS_PERMISSION", // (Jira Software only)
		"SERVICEDESK_AGENT",         // (Jira Service Desk only)
		"VIEW_DEV_TOOLS",            // (Jira Software only)
		"VIEW_READONLY_WORKFLOW",
		// Issue permissions
		"ASSIGNABLE_USER", "ASSIGN_ISSUES", "CLOSE_ISSUES", "CREATE_ISSUES", "DELETE_ISSUES", "EDIT_ISSUES", "LINK_ISSUES",
		"MODIFY_REPORTER", "MOVE_ISSUES", "RESOLVE_ISSUES", "SCHEDULE_ISSUES", "SET_ISSUE_SECURITY", "TRANSITION_ISSUES",
		// Voters and watchers permissions
		"MANAGE_WATCHERS", "VIEW_VOTERS_AND_WATCHERS",
		// Comments permissions
		"ADD_COMMENTS", "DELETE_ALL_COMMENTS", "DELETE_OWN_COMMENTS", "EDIT_ALL_COMMENTS", "EDIT_OWN_COMMENTS",
		// Attachments permissions
		"CREATE_ATTACHMENTS", "DELETE_ALL_ATTACHMENTS", "DELETE_OWN_ATTACHMENTS",
		// Time tracking permissions
		"DELETE_ALL_WORKLOGS", "DELETE_OWN_WORKLOGS", "EDIT_ALL_WORKLOGS", "EDIT_OWN_WORKLOGS", "WORK_ON_ISSUES",
	}
)

func NewJiraPermissionGrantResource() resource.Resource {
	return &jiraPermissionGrantResource{}
}

func (*jiraPermissionGrantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_permission_grant"
}

func (*jiraPermissionGrantResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Jira Permission Grant Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "(Forces new) The ID of the permission grant.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"permission_scheme_id": schema.StringAttribute{
				MarkdownDescription: "(Forces new) The ID of the permission scheme in which to create a new permission grant.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"holder": schema.SingleNestedAttribute{
				MarkdownDescription: "(Forces new) The user, group, field or role being granted the permission.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of permission holder. " +
							"Can be one of: `anyone`, `applicationRole`, `assignee`, `group`, `groupCustomField`, " +
							"`projectLead`, `projectRole`, `reporter`, `user` or `userCustomField`.",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(holder_types...),
						},
					},
					"parameter": schema.StringAttribute{
						MarkdownDescription: "The identifier associated with the `type` value that defines the holder of the permission.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringmodifiers.DefaultValue(""),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"permission": schema.StringAttribute{
				MarkdownDescription: "(Forces new) The permission to grant. Can be one of the built-in permissions or a custom permission added by an app.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(built_in_permissions...),
				},
			},
		},
	}
}

func (r *jiraPermissionGrantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jira.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jira.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.p.jira = client
}

func (*jiraPermissionGrantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ID, permission_scheme_id. Got: %q", req.ID))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Importing permission grant with import identifier: %+v", idParts))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("permission_scheme_id"), idParts[1])...)
}

func (r *jiraPermissionGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating permission grant resource")

	var plan jiraPermissionGrantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission grant plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v, Holder:%+v", plan, plan.Holder),
	})

	specialTypes := []string{"group", "projectRole", "user", "userCustomField"}
	for _, st := range specialTypes {
		if plan.Holder.Type.ValueString() == st {
			if plan.Holder.Parameter.ValueString() == "" {
				resp.Diagnostics.AddAttributeError(path.Root("holder").AtMapKey("parameter"),
					"Failed to provide a value for \"holder.parameter\" attribute",
					fmt.Sprintf("Value must be provided if \"holder.type\" is: %s", st),
				)
				return
			}
		}
	}

	schemeId, _ := strconv.Atoi(plan.PermissionSchemeID.ValueString())
	createPayload := &models.PermissionGrantPayloadScheme{
		Holder: &models.PermissionGrantHolderScheme{
			Type:      plan.Holder.Type.ValueString(),
			Parameter: plan.Holder.Parameter.ValueString(),
		},
		Permission: plan.Permission.ValueString(),
	}

	permissionGrant, res, err := r.p.jira.Permission.Scheme.Grant.Create(ctx, schemeId, createPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create permission grant, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created permission grant")

	plan.ID = types.StringValue(strconv.Itoa(permissionGrant.ID))

	tflog.Debug(ctx, "Storing permission grant into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v, Holder:%+v", plan, plan.Holder),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraPermissionGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading permission grant resource")

	var state jiraPermissionGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission grant from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v, Holder:%+v", state, state.Holder),
	})

	grantId, _ := strconv.Atoi(state.ID.ValueString())
	schemeId, _ := strconv.Atoi(state.PermissionSchemeID.ValueString())

	permissionGrant, res, err := r.p.jira.Permission.Scheme.Grant.Get(ctx, schemeId, grantId, []string{"all"})
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get permission grant, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved permission grant from API state")

	state.Holder = &jiraPermissionGrantHolderModel{
		Type:      types.StringValue(permissionGrant.Holder.Type),
		Parameter: types.StringValue(permissionGrant.Holder.Parameter),
	}
	state.Permission = types.StringValue(permissionGrant.Permission)

	tflog.Debug(ctx, "Storing permission grant into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v, Holder:%+v", state, state.Holder),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraPermissionGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The RequiresReplace plan modifier will trigger Terraform to destroy and recreate the resource
	// if any of the required attributes changes, i.e. permission_scheme_id, holder or permission
	tflog.Debug(ctx, "If the value of any required attribute changes, Terraform will destroy and recreate the resource")
}

func (r *jiraPermissionGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting permission grant resource")

	var state jiraPermissionGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	grantId, _ := strconv.Atoi(state.ID.ValueString())
	schemeId, _ := strconv.Atoi(state.PermissionSchemeID.ValueString())

	res, err := r.p.jira.Permission.Scheme.Grant.Delete(ctx, schemeId, grantId)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete permission grant, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted permission grant from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
