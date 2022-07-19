package atlassian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_validation"

	models "github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

var _ tfsdk.Resource = jiraIssueTypeSchemeResource{}
var _ tfsdk.ResourceWithImportState = jiraIssueTypeSchemeResource{}
var _ tfsdk.ResourceType = jiraIssueTypeSchemeResourceType{}

type jiraIssueTypeSchemeResource struct {
	p provider
}
type jiraIssueTypeSchemeResourceType struct{}
type jiraIssueTypeSchemeResourceData struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	DefaultIssueTypeId types.String `tfsdk:"default_issue_type_id"`
	IssueTypeIds       types.List   `tfsdk:"issue_type_ids"`
}

func (t jiraIssueTypeSchemeResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Jira Issue Type Scheme Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type scheme.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue type scheme. The name must be unique. The maximum length is 255 characters.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.StringLengthBetween(0, 255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the issue type scheme. The maximum length is 4000 characters.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.StringLengthBetween(0, 4000),
				},
			},
			"default_issue_type_id": {
				MarkdownDescription: "The ID of the default issue type of the issue type scheme. This ID must be included in issue_type_ids.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"issue_type_ids": {
				MarkdownDescription: "The list of issue types IDs of the issue type scheme. At least one standard issue type ID is required.",
				Required:            true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
		},
		Version: 1,
	}, nil
}

func (t jiraIssueTypeSchemeResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jiraIssueTypeSchemeResource{
		p: provider,
	}, diags
}

func (r jiraIssueTypeSchemeResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r jiraIssueTypeSchemeResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan jiraIssueTypeSchemeResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Description.Unknown {
		plan.Description = types.String{Value: ""}
	}

	if plan.DefaultIssueTypeId.Unknown {
		plan.DefaultIssueTypeId = types.String{Value: ""}
	}

	if plan.DefaultIssueTypeId.Value != "" {
		flag := false
		for _, v := range plan.IssueTypeIds.Elems {
			if v == plan.DefaultIssueTypeId {
				flag = true
			}
		}
		if !flag {
			resp.Diagnostics.AddError("User Error", "Value of default_issue_type_id must be included in issue_type_ids.")
			return
		}
	}

	issueTypeSchemePayload := new(models.IssueTypeSchemePayloadScheme)
	issueTypeSchemePayload.Name = plan.Name.Value
	issueTypeSchemePayload.Description = plan.Description.Value
	issueTypeSchemePayload.DefaultIssueTypeID = plan.DefaultIssueTypeId.Value
	diags = plan.IssueTypeIds.ElementsAs(ctx, &issueTypeSchemePayload.IssueTypeIds, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	returnedIssueTypeScheme, res, err := r.p.jira.Issue.Type.Scheme.Create(ctx, issueTypeSchemePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue type scheme, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}

	// Store new issue type scheme ID in state
	plan.ID = types.String{Value: returnedIssueTypeScheme.IssueTypeSchemeID}

	tflog.Debug(ctx, "created an issue type scheme", map[string]interface{}{
		"issue_type_scheme_id": returnedIssueTypeScheme.IssueTypeSchemeID,
	})

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r jiraIssueTypeSchemeResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state jiraIssueTypeSchemeResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	issueTypeSchemeID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Conversion failed: %s", err.Error()))
		return
	}

	// Get issue type scheme details
	returnedIssueTypeScheme, res, err := r.p.jira.Issue.Type.Scheme.Gets(ctx, []int{issueTypeSchemeID}, 0, 50)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read issue type scheme, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}

	// Get issue type scheme items
	returnedIssueTypeSchemeItems, res, err := r.p.jira.Issue.Type.Scheme.Items(ctx, []int{issueTypeSchemeID}, 0, 50)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type scheme items, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}
	ids := types.List{
		ElemType: types.StringType,
	}
	for _, elem := range returnedIssueTypeSchemeItems.Values {
		av := types.String{Value: elem.IssueTypeID}
		ids.Elems = append(ids.Elems, av)
	}

	state.Name = types.String{Value: returnedIssueTypeScheme.Values[0].Name}
	state.Description = types.String{Value: returnedIssueTypeScheme.Values[0].Description}
	state.DefaultIssueTypeId = types.String{Value: returnedIssueTypeScheme.Values[0].DefaultIssueTypeID}
	state.IssueTypeIds = ids

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r jiraIssueTypeSchemeResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan jiraIssueTypeSchemeResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state jiraIssueTypeSchemeResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	issueTypeSchemeID, _ := strconv.Atoi(state.ID.Value)

	issueTypeSchemePayload := new(models.IssueTypeSchemePayloadScheme)
	issueTypeSchemePayload.Name = plan.Name.Value
	issueTypeSchemePayload.Description = plan.Description.Value

	res, err := r.p.jira.Issue.Type.Scheme.Update(ctx, issueTypeSchemeID, issueTypeSchemePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue type scheme, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}

	// Validate that default_issue_type_id is included in issue_type_ids
	if plan.DefaultIssueTypeId.Value != "" {
		flag := false
		for _, v := range plan.IssueTypeIds.Elems {
			if v == plan.DefaultIssueTypeId {
				flag = true
			}
		}
		if !flag {
			resp.Diagnostics.AddError("User Error", "Value of default_issue_type_id must be included in issue_type_ids.")
			return
		}
	}

	// Validate that new issue type(s) need to be added to issue type scheme
	var ids []int
	var exists bool
	for _, p := range plan.IssueTypeIds.Elems {
		exists = false
		for _, s := range state.IssueTypeIds.Elems {
			if p == s {
				exists = true
			}
		}
		if !exists {
			new_id, _ := strconv.Atoi(p.String())
			ids = append(ids, new_id)
		}
	}

	// Add new issue type(s) to issue type scheme
	if len(ids) != 0 {
		res, err = r.p.jira.Issue.Type.Scheme.Append(ctx, issueTypeSchemeID, ids)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add issue types to issue type scheme, got error: %s\n%s", err.Error(), res.Bytes.String()))
			return
		}
	}

	var result = jiraIssueTypeSchemeResourceData{
		ID:                 types.String{Value: state.ID.Value},
		Name:               types.String{Value: plan.Name.Value},
		Description:        types.String{Value: plan.Description.Value},
		DefaultIssueTypeId: types.String{Value: plan.DefaultIssueTypeId.Value},
		IssueTypeIds:       plan.IssueTypeIds,
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r jiraIssueTypeSchemeResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state jiraIssueTypeSchemeResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	issueTypeSchemeID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to convert issue type scheme ID, got error: %s", err.Error()))
		return
	}

	res, err := r.p.jira.Issue.Type.Scheme.Delete(ctx, issueTypeSchemeID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue type scheme, got error: %s\n%s", err, res.Bytes.String()))
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}
