package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_plan_modification"
)

type (
	jiraScreenSchemeResource struct {
		p atlassianProvider
	}

	jiraScreenSchemeResourceModel struct {
		ID          types.String                `tfsdk:"id"`
		Name        types.String                `tfsdk:"name"`
		Description types.String                `tfsdk:"description"`
		Screens     *jiraScreenSchemeTypesModel `tfsdk:"screens"`
	}

	jiraScreenSchemeTypesModel struct {
		Create  types.Int64 `tfsdk:"create"`
		Default types.Int64 `tfsdk:"default"`
		View    types.Int64 `tfsdk:"view"`
		Edit    types.Int64 `tfsdk:"edit"`
	}
)

var (
	_ resource.Resource                = (*jiraScreenSchemeResource)(nil)
	_ resource.ResourceWithImportState = (*jiraScreenSchemeResource)(nil)
)

func NewJiraScreenSchemeResource() resource.Resource {
	return &jiraScreenSchemeResource{}
}

func (*jiraScreenSchemeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_screen_scheme"
}

func (*jiraScreenSchemeResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				PlanModifiers: tfsdk.AttributePlanModifiers{
					attribute_plan_modification.DefaultValue(types.StringValue("")),
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

func (r *jiraScreenSchemeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraScreenSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraScreenSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating screen scheme resource")

	var plan jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	createRequestPayload := models.ScreenSchemePayloadScheme{
		Screens: &models.ScreenTypesScheme{
			Create:  int(plan.Screens.Create.ValueInt64()),
			Default: int(plan.Screens.Default.ValueInt64()),
			View:    int(plan.Screens.View.ValueInt64()),
			Edit:    int(plan.Screens.Edit.ValueInt64()),
		},
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	screenScheme, res, err := r.p.jira.Screen.Scheme.Create(ctx, &createRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create screen scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created screen scheme")

	plan.ID = types.StringValue(strconv.Itoa(screenScheme.ID))

	tflog.Debug(ctx, "Storing screen scheme into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraScreenSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading screen scheme resource")

	var state jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	screenSchemeId, _ := strconv.Atoi(state.ID.ValueString())
	options := &models.ScreenSchemeParamsScheme{
		IDs: []int{screenSchemeId},
	}
	resScreenScheme, res, err := r.p.jira.Screen.Scheme.Gets(ctx, options, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get screen scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved screen scheme from API state")

	state.Name = types.String{Value: resScreenScheme.Values[0].Name}
	state.Description = types.String{Value: resScreenScheme.Values[0].Description}
	state.Screens = &jiraScreenSchemeTypesModel{
		Create:  types.Int64{Value: int64(resScreenScheme.Values[0].Screens.Create)},
		Default: types.Int64{Value: int64(resScreenScheme.Values[0].Screens.Default)},
		View:    types.Int64{Value: int64(resScreenScheme.Values[0].Screens.View)},
		Edit:    types.Int64{Value: int64(resScreenScheme.Values[0].Screens.Edit)},
	}
	tflog.Debug(ctx, "Storing screen scheme into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraScreenSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating screen scheme resource")

	var plan jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

	updateRequestPayload := models.ScreenSchemePayloadScheme{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Screens: &models.ScreenTypesScheme{
			Create:  int(plan.Screens.Create.ValueInt64()),
			Default: int(plan.Screens.Default.ValueInt64()),
			View:    int(plan.Screens.View.ValueInt64()),
			Edit:    int(plan.Screens.Edit.ValueInt64()),
		},
	}

	res, err := r.p.jira.Screen.Scheme.Update(ctx, state.ID.ValueString(), &updateRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update screen scheme, got error: %s\n%s", err, resBody))
	}
	tflog.Debug(ctx, "Updated screen scheme in API state")

	plan.ID = types.String{Value: state.ID.ValueString()}

	tflog.Debug(ctx, "Storing screen scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraScreenSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting screen scheme resource")

	var state jiraScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme from state")

	res, err := r.p.jira.Screen.Scheme.Delete(ctx, state.ID.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete screen scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted screen scheme from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
