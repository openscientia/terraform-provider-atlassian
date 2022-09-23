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
)

type (
	jiraIssueTypeSchemeResource struct {
		p atlassianProvider
	}

	jiraIssueTypeSchemeResourceModel struct {
		ID                 types.String `tfsdk:"id"`
		Name               types.String `tfsdk:"name"`
		Description        types.String `tfsdk:"description"`
		DefaultIssueTypeId types.String `tfsdk:"default_issue_type_id"`
		IssueTypeIds       types.List   `tfsdk:"issue_type_ids"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueTypeSchemeResource)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueTypeSchemeResource)(nil)
)

func NewJiraIssueTypeSchemeResource() resource.Resource {
	return &jiraIssueTypeSchemeResource{}
}

func (*jiraIssueTypeSchemeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_type_scheme"
}

func (*jiraIssueTypeSchemeResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
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
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the issue type scheme. The maximum length is 4000 characters.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(4000),
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
	}, nil
}

func (r *jiraIssueTypeSchemeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraIssueTypeSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraIssueTypeSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue type scheme resource")

	var plan jiraIssueTypeSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type scheme plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

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
	resp.Diagnostics.Append(plan.IssueTypeIds.ElementsAs(ctx, &issueTypeSchemePayload.IssueTypeIds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	returnedIssueTypeScheme, res, err := r.p.jira.Issue.Type.Scheme.Create(ctx, issueTypeSchemePayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue type scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue type scheme")

	plan.ID = types.String{Value: returnedIssueTypeScheme.IssueTypeSchemeID}

	tflog.Debug(ctx, "Storing issue type scheme into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueTypeSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type scheme resource")

	var state jiraIssueTypeSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type scheme from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	issueTypeSchemeID, _ := strconv.Atoi(state.ID.Value)

	returnedIssueTypeScheme, res, err := r.p.jira.Issue.Type.Scheme.Gets(ctx, []int{issueTypeSchemeID}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read issue type scheme, got error: %s\n%s", err, resBody))
		return
	}

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
	tflog.Debug(ctx, "Retrieved issue type scheme from API state")

	state.Name = types.String{Value: returnedIssueTypeScheme.Values[0].Name}
	state.Description = types.String{Value: returnedIssueTypeScheme.Values[0].Description}
	state.DefaultIssueTypeId = types.String{Value: returnedIssueTypeScheme.Values[0].DefaultIssueTypeID}
	state.IssueTypeIds = ids

	tflog.Debug(ctx, "Storing issue type scheme into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueTypeSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue type scheme resource")

	var plan jiraIssueTypeSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type scheme plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueTypeSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type scheme from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

	issueTypeSchemeID, _ := strconv.Atoi(state.ID.Value)

	issueTypeSchemePayload := new(models.IssueTypeSchemePayloadScheme)
	issueTypeSchemePayload.Name = plan.Name.Value
	issueTypeSchemePayload.Description = plan.Description.Value

	res, err := r.p.jira.Issue.Type.Scheme.Update(ctx, issueTypeSchemeID, issueTypeSchemePayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue type scheme, got error: %s\n%s", err, resBody))
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
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add issue types to issue type scheme, got error: %s\n%s", err, resBody))
			return
		}
	}
	tflog.Debug(ctx, "Updated issue type scheme in API state")

	var result = jiraIssueTypeSchemeResourceModel{
		ID:                 types.String{Value: state.ID.Value},
		Name:               types.String{Value: plan.Name.Value},
		Description:        types.String{Value: plan.Description.Value},
		DefaultIssueTypeId: types.String{Value: plan.DefaultIssueTypeId.Value},
		IssueTypeIds:       plan.IssueTypeIds,
	}

	tflog.Debug(ctx, "Storing issue type scheme into the state", map[string]interface{}{
		"updateNewState": fmt.Sprintf("%+v", result),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *jiraIssueTypeSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue type scheme resource")

	var state jiraIssueTypeSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type scheme from state")

	issueTypeSchemeID, _ := strconv.Atoi(state.ID.Value)

	res, err := r.p.jira.Issue.Type.Scheme.Delete(ctx, issueTypeSchemeID)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue type scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted issue type scheme from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
