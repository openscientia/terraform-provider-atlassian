package atlassian

import (
	"context"
	"fmt"
	"strconv"

	models "github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_plan_modification"
)

type (
	jiraScreenSchemeResource struct {
		p provider
	}

	jiraScreenSchemeResourceType struct{}

	jiraScreenSchemeResourceModel struct {
		ID          types.String               `tfsdk:"id"`
		Name        types.String               `tfsdk:"name"`
		Description types.String               `tfsdk:"description"`
		Screens     jiraScreenSchemeTypesModel `tfsdk:"screens"`
	}
	jiraScreenSchemeTypesModel struct {
		Create  types.Int64 `tfsdk:"create"`
		Default types.Int64 `tfsdk:"default"`
		View    types.Int64 `tfsdk:"view"`
		Edit    types.Int64 `tfsdk:"edit"`
	}
)

var (
	_ tfsdk.Resource                = (*jiraScreenSchemeResource)(nil)
	_ tfsdk.ResourceType            = (*jiraScreenSchemeResourceType)(nil)
	_ tfsdk.ResourceWithImportState = (*jiraScreenSchemeResource)(nil)
)

func (*jiraScreenSchemeResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Screen Scheme Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the screen scheme.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the screen scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the screen scheme. " +
					"The maximum length is 255 characters.",
				Optional: true,
				Computed: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(255),
				},
			},
			"screens": {
				MarkdownDescription: "The IDs of the screens for the screen types of the screen scheme. " +
					"Only screens used in classic projects are accepted.",
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"create": {
							MarkdownDescription: "The ID of the create screen.",
							Optional:            true,
							Computed:            true,
							Type:                types.Int64Type,
							PlanModifiers: tfsdk.AttributePlanModifiers{
								attribute_plan_modification.DefaultValue(types.Int64{Value: 0}),
							},
						},
						"default": {
							MarkdownDescription: "The ID of the default screen. Required when creating a screen scheme.",
							Required:            true,
							Type:                types.Int64Type,
						},
						"view": {
							MarkdownDescription: "The ID of the view screen.",
							Optional:            true,
							Computed:            true,
							Type:                types.Int64Type,
							PlanModifiers: tfsdk.AttributePlanModifiers{
								attribute_plan_modification.DefaultValue(types.Int64{Value: 0}),
							},
						},
						"edit": {
							MarkdownDescription: "The ID of the edit screen.",
							Optional:            true,
							Computed:            true,
							Type:                types.Int64Type,
							PlanModifiers: tfsdk.AttributePlanModifiers{
								attribute_plan_modification.DefaultValue(types.Int64{Value: 0}),
							},
						},
					},
				),
			},
		},
	}, nil
}

func (r *jiraScreenSchemeResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraScreenSchemeResource{
		p: provider,
	}, diags
}

func (r *jiraScreenSchemeResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	// Must initialise "singleNestedAttributes" to avoid "unhandled null values" error when calling (tfsdk.Plan).Get
	screens := jiraScreenSchemeTypesModel{}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("screens"), screens)...)
}

func (r *jiraScreenSchemeResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating screen scheme resource")

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
	}

	var plan jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme configuration", map[string]interface{}{
		"screenSchemeConfig": fmt.Sprintf("%+v", plan),
	})

	if plan.Description.IsUnknown() || plan.Description.IsNull() {
		plan.Description = types.String{Value: ""}
	}

	createRequestPayload := models.ScreenSchemePayloadScheme{
		Screens: &models.ScreenTypesScheme{
			Create:  int(plan.Screens.Create.Value),
			Default: int(plan.Screens.Default.Value),
			View:    int(plan.Screens.View.Value),
			Edit:    int(plan.Screens.Edit.Value),
		},
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
	}
	tflog.Debug(ctx, "Generated request payload", map[string]interface{}{
		"screenSchemeReq": fmt.Sprintf("%+v", createRequestPayload),
	})
	newScreenScheme, res, err := r.p.jira.Screen.Scheme.Create(ctx, &createRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create screen scheme, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Created screen scheme", map[string]interface{}{
		"screenScheme": newScreenScheme.ID,
	})

	plan.ID = types.String{Value: strconv.Itoa(newScreenScheme.ID)}

	tflog.Debug(ctx, "Storing screen scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraScreenSchemeResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Reading screen scheme resource")

	var state jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme from state", map[string]interface{}{
		"screenSchemeState": fmt.Sprintf("%+v", state),
	})

	screenSchemeId, _ := strconv.Atoi(state.ID.Value)
	resScreenScheme, res, err := r.p.jira.Screen.Scheme.Gets(ctx, []int{screenSchemeId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get screen scheme, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved screen scheme from API state")

	state.Name = types.String{Value: resScreenScheme.Values[0].Name}
	state.Description = types.String{Value: resScreenScheme.Values[0].Description}
	state.Screens = jiraScreenSchemeTypesModel{
		Create:  types.Int64{Value: int64(resScreenScheme.Values[0].Screens.Create)},
		Default: types.Int64{Value: int64(resScreenScheme.Values[0].Screens.Default)},
		View:    types.Int64{Value: int64(resScreenScheme.Values[0].Screens.View)},
		Edit:    types.Int64{Value: int64(resScreenScheme.Values[0].Screens.Edit)},
	}
	tflog.Debug(ctx, "Storing screen scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraScreenSchemeResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating screen scheme")

	var plan jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded screen scheme configuration", map[string]interface{}{
		"screenSchemeConfig": fmt.Sprintf("%+v", plan),
	})

	var state jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded screen scheme from state", map[string]interface{}{
		"screenSchemestate": fmt.Sprintf("%+v", state),
	})

	updateRequestPayload := models.ScreenSchemePayloadScheme{
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
		Screens: &models.ScreenTypesScheme{
			Create:  int(plan.Screens.Create.Value),
			Default: int(plan.Screens.Default.Value),
			View:    int(plan.Screens.View.Value),
			Edit:    int(plan.Screens.Edit.Value),
		},
	}
	tflog.Debug(ctx, "Generated request payload", map[string]interface{}{
		"screenSchemeReq": fmt.Sprintf("%+v", updateRequestPayload),
	})
	res, err := r.p.jira.Screen.Scheme.Update(ctx, state.ID.Value, &updateRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update screen scheme, got error: %s\n%s", err.Error(), resBody))
	}
	tflog.Debug(ctx, "Updated screen scheme in API state")

	plan.ID = types.String{Value: state.ID.Value}

	tflog.Debug(ctx, "Storing screen scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraScreenSchemeResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Deleting screen scheme resource")

	var state jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme from state")

	res, err := r.p.jira.Screen.Scheme.Delete(ctx, state.ID.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete screen scheme, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Removed screen scheme from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
