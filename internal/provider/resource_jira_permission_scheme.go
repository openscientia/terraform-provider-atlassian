package atlassian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_plan_modification"
)

type (
	jiraPermissionSchemeResource struct {
		p atlassianProvider
	}

	jiraPermissionSchemeResourceType struct{}

	jiraPermissionSchemeResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Self        types.String `tfsdk:"self"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ resource.Resource                = (*jiraPermissionSchemeResource)(nil)
	_ provider.ResourceType            = (*jiraPermissionSchemeResourceType)(nil)
	_ resource.ResourceWithImportState = (*jiraPermissionSchemeResource)(nil)
)

func (*jiraPermissionSchemeResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Permission Scheme Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the permission scheme.",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"self": {
				MarkdownDescription: "The URL of the permission scheme.",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "The name of the permission scheme. " +
					"The name must be unique. The maximum length is 255 characters.",
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the permission scheme.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					attribute_plan_modification.DefaultValue(types.String{Value: ""}),
				},
			},
		},
	}, nil
}

func (r *jiraPermissionSchemeResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraPermissionSchemeResource{
		p: provider,
	}, diags
}

func (r *jiraPermissionSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraPermissionSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating permission scheme resource")

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
	}

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
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
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

	plan.ID = types.String{Value: strconv.Itoa(permissionScheme.ID)}
	plan.Self = types.String{Value: permissionScheme.Self}

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

	schemeId, _ := strconv.Atoi(state.ID.Value)

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

	state.Self = types.String{Value: permissionScheme.Self}
	state.Name = types.String{Value: permissionScheme.Name}
	state.Description = types.String{Value: permissionScheme.Description}

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

	schemeId, _ := strconv.Atoi(state.ID.Value)

	updatePayload := &models.PermissionSchemeScheme{
		ID:          schemeId,
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
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

	schemeId, _ := strconv.Atoi(state.ID.Value)

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
