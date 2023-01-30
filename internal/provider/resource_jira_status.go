package atlassian

import (
	"context"
	"fmt"

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
	jiraStatusResource struct {
		p atlassianProvider
	}

	jiraStatusResourceModel struct {
		ID             types.String          `tfsdk:"id"`
		Name           types.String          `tfsdk:"name"`
		StatusCategory types.String          `tfsdk:"status_category"`
		Description    types.String          `tfsdk:"description"`
		StatusScope    *jiraStatusScopeModel `tfsdk:"status_scope"`
	}
	jiraStatusScopeModel struct {
		Type types.String `tfsdk:"type"`
		Id   types.String `tfsdk:"id"`
	}
)

var (
	_ resource.Resource                = (*jiraStatusResource)(nil)
	_ resource.ResourceWithImportState = (*jiraStatusResource)(nil)
)

func NewJiraStatusResource() resource.Resource {
	return &jiraStatusResource{}
}

func (*jiraStatusResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_status"
}

func (*jiraStatusResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Jira Status Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the status.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the status.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"status_category": schema.StringAttribute{
				MarkdownDescription: "The category of the status. Can be one of: `TODO`, `IN_PROGRESS`, `DONE`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("TODO", "IN_PROGRESS", "DONE"),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the status.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					// attribute must have value to avoid http error when creating status via the api endpoint:
					// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-status/#api-rest-api-3-statuses-post
					stringmodifiers.DefaultValue(" "),
				},
			},
			"status_scope": schema.SingleNestedAttribute{
				MarkdownDescription: "The scope of the status.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "(Forces new) The scope of the status. `GLOBAL` for company-managed projects and `PROJECT` for team-managed projects.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("GLOBAL", "PROJECT"),
						},
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "(Forces new) The ID of a team-managed project. Only use when `status_scope.type` is `PROJECT`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringmodifiers.DefaultValue(""),
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	}
}

func (r *jiraStatusResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraStatusResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraStatusResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating status resource")

	var plan jiraStatusResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded status plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v, %+v", plan, plan.StatusScope),
	})

	if plan.StatusScope.Type.ValueString() == "GLOBAL" {
		if plan.StatusScope.Id.ValueString() != "" {
			resp.Diagnostics.AddAttributeError(path.Root("status_scope").AtMapKey("id"),
				"\"GLOBAL\" scope types must not have a value for \"status_scope.id\" attribute",
				fmt.Sprintf("A value must not be provided if \"status_scope.type\" is: %s", plan.StatusScope.Type.ValueString()))
			return
		}
	}

	if plan.StatusScope.Type.ValueString() == "PROJECT" {
		if plan.StatusScope.Id.ValueString() == "" || plan.StatusScope.Id.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("status_scope").AtMapKey("id"),
				"Failed to provide value for \"status_scope.id\" attribute",
				fmt.Sprintf("A value must be provided if \"status_scope.type\" is: %s", plan.StatusScope.Type.ValueString()))
			return
		}
	}

	payload := &models.WorkflowStatusPayloadScheme{}
	payload.Statuses = []*models.WorkflowStatusNodeScheme{
		{
			Name:           plan.Name.ValueString(),
			StatusCategory: plan.StatusCategory.ValueString(),
			Description:    plan.Description.ValueString(),
		},
	}
	if plan.StatusScope.Id.IsNull() || plan.StatusScope.Id.IsUnknown() || plan.StatusScope.Id.ValueString() == "" {
		payload.Scope = &models.WorkflowStatusScopeScheme{
			Type:    plan.StatusScope.Type.ValueString(),
			Project: nil,
		}
	} else {
		payload.Scope = &models.WorkflowStatusScopeScheme{
			Type: plan.StatusScope.Type.ValueString(),
			Project: &models.WorkflowStatusProjectScheme{
				ID: plan.StatusScope.Id.ValueString(),
			},
		}
	}

	status, res, err := r.p.jira.Workflow.Status.Create(ctx, payload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created status in API state")

	plan.ID = types.StringValue(status[0].ID)

	tflog.Debug(ctx, "Storing status into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraStatusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading status resource")

	var state jiraStatusResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded status from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	status, res, err := r.p.jira.Workflow.Status.Gets(ctx, []string{state.ID.ValueString()}, []string{})
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get status, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved status from API state", map[string]interface{}{
		"status": fmt.Sprintf("%v", &status[0].Scope.Type),
	})

	state.Name = types.StringValue(status[0].Name)
	state.Description = types.StringValue(status[0].Description)
	state.StatusCategory = types.StringValue(status[0].StatusCategory)
	state.StatusScope = &jiraStatusScopeModel{
		Type: types.StringValue(status[0].Scope.Type),
	}
	if status[0].Scope.Project != nil {
		state.StatusScope.Id = types.StringValue(status[0].Scope.Project.ID)
	} else {
		state.StatusScope.Id = types.StringValue("")
	}

	tflog.Debug(ctx, "Storing status into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraStatusResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating status resource")

	var plan jiraStatusResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded status plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraStatusResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded status from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

	payload := &models.WorkflowStatusPayloadScheme{
		Statuses: []*models.WorkflowStatusNodeScheme{
			{
				ID:             state.ID.ValueString(),
				Name:           plan.Name.ValueString(),
				Description:    plan.Description.ValueString(),
				StatusCategory: plan.StatusCategory.ValueString(),
			},
		},
	}

	res, err := r.p.jira.Workflow.Status.Update(ctx, payload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Updated status in API state")

	tflog.Debug(ctx, "Storing status into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraStatusResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting status resource")

	var state jiraStatusResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded status from state")

	res, err := r.p.jira.Workflow.Status.Delete(ctx, []string{state.ID.ValueString()})
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted status from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
