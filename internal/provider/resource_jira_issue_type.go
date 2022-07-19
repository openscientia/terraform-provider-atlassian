package atlassian

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_validation"

	models "github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

var _ tfsdk.Resource = jiraIssueTypeResource{}
var _ tfsdk.ResourceType = jiraIssueTypeResourceType{}
var _ tfsdk.ResourceWithImportState = jiraIssueTypeResource{}

type jiraIssueTypeResourceType struct{}

type jiraIssueTypeResource struct {
	p provider
}

type jiraIssueTypeResourceData struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	HierarchyLevel types.Int64  `tfsdk:"hierarchy_level"`
	AvatarId       types.Int64  `tfsdk:"avatar_id"`
}

func (t jiraIssueTypeResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Jira Issue Type Resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue type. The maximum length is 60 characters.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.StringLengthBetween(0, 60),
				},
			},
			"description": {
				MarkdownDescription: "The description of the issue type.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "The type of the issue type. Can be either `standard` or `sub-task`.",
				DeprecationMessage:  "Use hierarchy_level instead.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.StringValues([]string{"standard", "sub-task"}),
				},
			},
			"hierarchy_level": {
				MarkdownDescription: "The hierarchy level of the issue type. Can be either `0` or `-1`.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.IntValues([]int{0, -1}),
				},
			},
			"avatar_id": {
				MarkdownDescription: "The ID of the issue type's avatar.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
		},
		Version: 1,
	}, nil
}

func (t jiraIssueTypeResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jiraIssueTypeResource{
		p: provider,
	}, diags
}

func (r jiraIssueTypeResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r jiraIssueTypeResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan jiraIssueTypeResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Type.Unknown && !plan.HierarchyLevel.Unknown {
		resp.Diagnostics.AddError("User Error", "Cannot use attributes `type` and `hierarchy_level` together.")
		return
	}

	if plan.Description.IsUnknown() {
		plan.Description = types.String{Value: ""}
	}

	if plan.Type.Unknown && plan.HierarchyLevel.Unknown {
		plan.Type = types.String{Value: "standard"}
		plan.HierarchyLevel = types.Int64{Value: 0}
	} else if plan.Type.Unknown && !plan.HierarchyLevel.Unknown {
		if plan.HierarchyLevel.Value == 0 {
			plan.Type = types.String{Value: "standard"}
		} else {
			plan.Type = types.String{Value: "sub-task"}
		}
	} else if !plan.Type.Unknown && plan.HierarchyLevel.Unknown {
		if plan.Type.Value == "standard" {
			plan.HierarchyLevel = types.Int64{Value: 0}
		} else {
			plan.HierarchyLevel = types.Int64{Value: -1}
		}
	}

	issueTypePayload := new(models.IssueTypePayloadScheme)
	issueTypePayload.Name = plan.Name.Value
	issueTypePayload.Description = plan.Description.Value
	issueTypePayload.HierarchyLevel = int(plan.HierarchyLevel.Value)

	// Create new issue type
	returnedIssueType, res, err := r.p.jira.Issue.Type.Create(ctx, issueTypePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}

	// Store new issue type ID in state
	plan.ID = types.String{Value: returnedIssueType.ID}

	// Apply chosen avatar image for new issue type
	if !plan.AvatarId.Unknown {
		issueTypePayload := new(models.IssueTypePayloadScheme)
		issueTypePayload.Name = plan.Name.Value
		issueTypePayload.Description = plan.Description.Value
		issueTypePayload.AvatarID = int(plan.AvatarId.Value)

		returnedIssueType, res, err := r.p.jira.Issue.Type.Update(ctx, returnedIssueType.ID, issueTypePayload)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
			return
		}
		plan.AvatarId = types.Int64{Value: int64(returnedIssueType.AvatarID)}

	} else {
		plan.AvatarId = types.Int64{Value: int64(returnedIssueType.AvatarID)}
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r jiraIssueTypeResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state jiraIssueTypeResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	issueTypeID := state.ID.Value

	returnedIssueType, res, err := r.p.jira.Issue.Type.Get(ctx, issueTypeID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}

	state.Name = types.String{Value: returnedIssueType.Name}
	state.Description = types.String{Value: returnedIssueType.Description}
	if returnedIssueType.HierarchyLevel == 0 {
		state.Type = types.String{Value: "standard"}
	} else {
		state.Type = types.String{Value: "sub-task"}
	}
	state.HierarchyLevel = types.Int64{Value: int64(returnedIssueType.HierarchyLevel)}
	state.AvatarId = types.Int64{Value: int64(returnedIssueType.AvatarID)}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r jiraIssueTypeResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan jiraIssueTypeResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state jiraIssueTypeResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	issueTypeID := state.ID.Value

	issueTypePayload := new(models.IssueTypePayloadScheme)
	issueTypePayload.Name = plan.Name.Value
	issueTypePayload.Description = plan.Description.Value
	issueTypePayload.AvatarID = int(plan.AvatarId.Value)

	returnedIssueType, res, err := r.p.jira.Issue.Type.Update(ctx, issueTypeID, issueTypePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}

	var result = jiraIssueTypeResourceData{
		ID:             types.String{Value: returnedIssueType.ID},
		Description:    types.String{Value: returnedIssueType.Description},
		Name:           types.String{Value: returnedIssueType.Name},
		Type:           types.String{Value: state.Type.Value},
		AvatarId:       types.Int64{Value: int64(returnedIssueType.AvatarID)},
		HierarchyLevel: types.Int64{Value: int64(returnedIssueType.HierarchyLevel)},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r jiraIssueTypeResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state jiraIssueTypeResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	issueTypeID := state.ID.Value

	res, err := r.p.jira.Issue.Type.Delete(ctx, issueTypeID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue type, got error: %s\n%s", err, res.Bytes.String()))
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}
