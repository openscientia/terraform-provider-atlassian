package atlassian

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	jiraIssueTypeResource struct {
		p atlassianProvider
	}

	jiraIssueTypeResourceModel struct {
		ID             types.String `tfsdk:"id"`
		Name           types.String `tfsdk:"name"`
		Description    types.String `tfsdk:"description"`
		Type           types.String `tfsdk:"type"`
		HierarchyLevel types.Int64  `tfsdk:"hierarchy_level"`
		AvatarId       types.Int64  `tfsdk:"avatar_id"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueTypeResource)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueTypeResource)(nil)
)

func NewJiraIssueTypeResource() resource.Resource {
	return &jiraIssueTypeResource{}
}

func (*jiraIssueTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_type"
}

func (*jiraIssueTypeResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Type Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue type. The maximum length is 60 characters.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(60),
				},
			},
			"description": {
				MarkdownDescription: "The description of the issue type.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					attribute_plan_modification.DefaultValue(types.StringValue("")),
				},
			},
			"type": {
				MarkdownDescription: "The type of the issue type. Can be either `standard` or `sub-task`.",
				DeprecationMessage:  "Use hierarchy_level instead.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("standard", "sub-task"),
				},
			},
			"hierarchy_level": {
				MarkdownDescription: "The hierarchy level of the issue type. Can be either `0` or `-1`.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					int64validator.OneOf(0, -1),
				},
			},
			"avatar_id": {
				MarkdownDescription: "The ID of the issue type's avatar.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
		},
	}, nil
}

func (r *jiraIssueTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraIssueTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraIssueTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue type resource")

	var plan jiraIssueTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	if !plan.Type.IsUnknown() && !plan.HierarchyLevel.IsUnknown() {
		resp.Diagnostics.AddError("User Error", "Cannot use attributes `type` and `hierarchy_level` together.")
		return
	}

	if plan.Type.IsUnknown() && plan.HierarchyLevel.IsUnknown() {
		plan.Type = types.StringValue("standard")
		plan.HierarchyLevel = types.Int64Value(0)
	} else if plan.Type.IsUnknown() && !plan.HierarchyLevel.IsUnknown() {
		if plan.HierarchyLevel.ValueInt64() == 0 {
			plan.Type = types.StringValue("standard")
		} else {
			plan.Type = types.StringValue("sub-task")
		}
	} else if !plan.Type.IsUnknown() && plan.HierarchyLevel.IsUnknown() {
		if plan.Type.ValueString() == "standard" {
			plan.HierarchyLevel = types.Int64Value(0)
		} else {
			plan.HierarchyLevel = types.Int64Value(-1)
		}
	}

	issueTypePayload := new(models.IssueTypePayloadScheme)
	issueTypePayload.Name = plan.Name.ValueString()
	issueTypePayload.Description = plan.Description.ValueString()
	issueTypePayload.HierarchyLevel = int(plan.HierarchyLevel.ValueInt64())

	returnedIssueType, res, err := r.p.jira.Issue.Type.Create(ctx, issueTypePayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue type, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue type")

	plan.ID = types.StringValue(returnedIssueType.ID)

	if !plan.AvatarId.IsUnknown() {
		issueTypePayload := new(models.IssueTypePayloadScheme)
		issueTypePayload.Name = plan.Name.ValueString()
		issueTypePayload.Description = plan.Description.ValueString()
		issueTypePayload.AvatarID = int(plan.AvatarId.ValueInt64())

		returnedIssueType, res, err := r.p.jira.Issue.Type.Update(ctx, returnedIssueType.ID, issueTypePayload)
		if err != nil {
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue type, got error: %s\n%s", err, resBody))
			return
		}
		plan.AvatarId = types.Int64Value(int64(returnedIssueType.AvatarID))
	} else {
		plan.AvatarId = types.Int64Value(int64(returnedIssueType.AvatarID))
	}

	tflog.Debug(ctx, "Storing issue type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type resource")

	var state jiraIssueTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	issueTypeID := state.ID.ValueString()

	returnedIssueType, res, err := r.p.jira.Issue.Type.Get(ctx, issueTypeID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}
	tflog.Debug(ctx, "Retrieved issue type from API state")

	state.Name = types.StringValue(returnedIssueType.Name)
	state.Description = types.StringValue(returnedIssueType.Description)
	if returnedIssueType.HierarchyLevel == 0 {
		state.Type = types.StringValue("standard")
	} else {
		state.Type = types.StringValue("sub-task")
	}
	state.HierarchyLevel = types.Int64Value(int64(returnedIssueType.HierarchyLevel))
	state.AvatarId = types.Int64Value(int64(returnedIssueType.AvatarID))

	tflog.Debug(ctx, "Storing issue type into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue type resource")

	var plan jiraIssueTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

	issueTypeID := state.ID.ValueString()

	issueTypePayload := new(models.IssueTypePayloadScheme)
	issueTypePayload.Name = plan.Name.ValueString()
	issueTypePayload.Description = plan.Description.ValueString()
	issueTypePayload.AvatarID = int(plan.AvatarId.ValueInt64())

	returnedIssueType, res, err := r.p.jira.Issue.Type.Update(ctx, issueTypeID, issueTypePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}
	tflog.Debug(ctx, "Updated issue type in API state")

	var result = jiraIssueTypeResourceModel{
		ID:             types.StringValue(returnedIssueType.ID),
		Description:    types.StringValue(returnedIssueType.Description),
		Name:           types.StringValue(returnedIssueType.Name),
		Type:           types.StringValue(state.Type.ValueString()),
		AvatarId:       types.Int64Value(int64(returnedIssueType.AvatarID)),
		HierarchyLevel: types.Int64Value(int64(returnedIssueType.HierarchyLevel)),
	}

	tflog.Debug(ctx, "Storing issue type into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *jiraIssueTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue type resource")

	var state jiraIssueTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type from state")

	res, err := r.p.jira.Issue.Type.Delete(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue type, got error: %s\n%s", err, res.Bytes.String()))
		return
	}
	tflog.Debug(ctx, "Deleted issue type from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
