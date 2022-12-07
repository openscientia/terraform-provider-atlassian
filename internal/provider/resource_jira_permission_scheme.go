package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/planmodifiers/stringmodifiers"
)

type (
	jiraPermissionSchemeResource struct {
		p atlassianProvider
	}

	jiraPermissionSchemeResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Self        types.String `tfsdk:"self"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ resource.Resource                = (*jiraPermissionSchemeResource)(nil)
	_ resource.ResourceWithImportState = (*jiraPermissionSchemeResource)(nil)
)

func NewJiraPermissionSchemeResource() resource.Resource {
	return &jiraPermissionSchemeResource{}
}

func (*jiraPermissionSchemeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_permission_scheme"
}

func (*jiraPermissionSchemeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Jira Permission Scheme Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the permission scheme.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"self": schema.StringAttribute{
				MarkdownDescription: "The URL of the permission scheme.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the permission scheme. " +
					"The name must be unique. The maximum length is 255 characters.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the permission scheme.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringmodifiers.DefaultValue(""),
				},
			},
		},
	}
}

func (r *jiraPermissionSchemeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraPermissionSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraPermissionSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating permission scheme resource")

	var plan jiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission scheme plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	createPayload := &models.PermissionSchemeScheme{
		Expand:      "all",
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	permissionScheme, res, err := r.p.jira.Permission.Scheme.Create(ctx, createPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create permission scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created permission scheme in API state")

	plan.ID = types.StringValue(strconv.Itoa(permissionScheme.ID))
	plan.Self = types.StringValue(permissionScheme.Self)

	tflog.Debug(ctx, "Storing permission scheme into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraPermissionSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading permission scheme resource")

	var state jiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission scheme from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	schemeId, _ := strconv.Atoi(state.ID.ValueString())

	permissionScheme, res, err := r.p.jira.Permission.Scheme.Get(ctx, schemeId, []string{""})
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get permission scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved permission scheme from API state")

	state.Self = types.StringValue(permissionScheme.Self)
	state.Name = types.StringValue(permissionScheme.Name)
	state.Description = types.StringValue(permissionScheme.Description)

	tflog.Debug(ctx, "Storing permission scheme into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraPermissionSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating permission scheme resource")

	var plan jiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission scheme plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission scheme from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

	schemeId, _ := strconv.Atoi(state.ID.ValueString())

	updatePayload := &models.PermissionSchemeScheme{
		ID:          schemeId,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	_, res, err := r.p.jira.Permission.Scheme.Update(ctx, schemeId, updatePayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update permission scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Updated permission scheme in API state")

	tflog.Debug(ctx, "Storing permission scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraPermissionSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting permission scheme resource")

	var state jiraPermissionSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission scheme from state")

	schemeId, _ := strconv.Atoi(state.ID.ValueString())

	res, err := r.p.jira.Permission.Scheme.Delete(ctx, schemeId)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete permission scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted permission scheme from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
