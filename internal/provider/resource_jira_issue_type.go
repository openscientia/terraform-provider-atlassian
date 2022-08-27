package atlassian

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = jiraIssueTypeResource{}
var _ provider.ResourceType = jiraIssueTypeResourceType{}
var _ resource.ResourceWithImportState = jiraIssueTypeResource{}

type jiraIssueTypeResourceType struct{}

type jiraIssueTypeResource struct {
	p atlassianProvider
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
		Version: 1,
	}, nil
}

func (t jiraIssueTypeResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jiraIssueTypeResource{
		p: provider,
	}, diags
}

func (r jiraIssueTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r jiraIssueTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue type resource")

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan jiraIssueTypeResourceData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

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

	plan.ID = types.String{Value: returnedIssueType.ID}

	if !plan.AvatarId.Unknown {
		issueTypePayload := new(models.IssueTypePayloadScheme)
		issueTypePayload.Name = plan.Name.Value
		issueTypePayload.Description = plan.Description.Value
		issueTypePayload.AvatarID = int(plan.AvatarId.Value)

		returnedIssueType, res, err := r.p.jira.Issue.Type.Update(ctx, returnedIssueType.ID, issueTypePayload)
		if err != nil {
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue type, got error: %s\n%s", err, resBody))
			return
		}
		plan.AvatarId = types.Int64{Value: int64(returnedIssueType.AvatarID)}
	} else {
		plan.AvatarId = types.Int64{Value: int64(returnedIssueType.AvatarID)}
	}

	tflog.Debug(ctx, "Storing issue type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r jiraIssueTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type resource")

	var state jiraIssueTypeResourceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	issueTypeID := state.ID.Value

	returnedIssueType, res, err := r.p.jira.Issue.Type.Get(ctx, issueTypeID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}
	tflog.Debug(ctx, "Retrieved issue type from API state")

	state.Name = types.String{Value: returnedIssueType.Name}
	state.Description = types.String{Value: returnedIssueType.Description}
	if returnedIssueType.HierarchyLevel == 0 {
		state.Type = types.String{Value: "standard"}
	} else {
		state.Type = types.String{Value: "sub-task"}
	}
	state.HierarchyLevel = types.Int64{Value: int64(returnedIssueType.HierarchyLevel)}
	state.AvatarId = types.Int64{Value: int64(returnedIssueType.AvatarID)}

	tflog.Debug(ctx, "Storing issue type into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r jiraIssueTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue type resource")

	var plan jiraIssueTypeResourceData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueTypeResourceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

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
	tflog.Debug(ctx, "Updated issue type in API state")

	var result = jiraIssueTypeResourceData{
		ID:             types.String{Value: returnedIssueType.ID},
		Description:    types.String{Value: returnedIssueType.Description},
		Name:           types.String{Value: returnedIssueType.Name},
		Type:           types.String{Value: state.Type.Value},
		AvatarId:       types.Int64{Value: int64(returnedIssueType.AvatarID)},
		HierarchyLevel: types.Int64{Value: int64(returnedIssueType.HierarchyLevel)},
	}

	tflog.Debug(ctx, "Storing issue type into the state", map[string]interface{}{
		"updateNewState": fmt.Sprintf("%+v", result),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r jiraIssueTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue type resource")

	var state jiraIssueTypeResourceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type from state")

	res, err := r.p.jira.Issue.Type.Delete(ctx, state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue type, got error: %s\n%s", err, res.Bytes.String()))
		return
	}
	tflog.Debug(ctx, "Deleted issue type from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
