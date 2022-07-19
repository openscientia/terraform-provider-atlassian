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
)

type (
	jiraIssueScreenResource struct {
		p provider
	}

	jiraIssueScreenResourceType struct{}

	jiraIssueScreenResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ tfsdk.Resource                = jiraIssueScreenResource{}
	_ tfsdk.ResourceType            = jiraIssueScreenResourceType{}
	_ tfsdk.ResourceWithImportState = jiraIssueScreenResource{}
)

func (jiraIssueScreenResourceType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Screen Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue screen.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the screen." +
					"The name must be unique." +
					"The maximum length is 255 characters.",
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.StringLengthBetween(0, 255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the screen." +
					"The maximum length is 255 characters.",
				Optional: true,
				Computed: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.StringLengthBetween(0, 255),
				},
			},
		},
	}, nil
}

func (jiraIssueScreenResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jiraIssueScreenResource{
		p: provider,
	}, diags

}

func (jiraIssueScreenResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r jiraIssueScreenResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating issue screen resource")

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan jiraIssueScreenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded issue screen configuration", map[string]interface{}{
		"issueScreenConfig": fmt.Sprintf("%+v", plan),
	})

	if plan.Description.Unknown {
		plan.Description = types.String{Value: ""}
	}

	newIssueScreen, res, err := r.p.jira.Screen.Create(ctx, plan.Name.Value, plan.Description.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue screen, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Created issue screen", map[string]interface{}{
		"issueScreen": newIssueScreen.ID,
	})

	tflog.Debug(ctx, "Storing issue screen info into the state")
	plan.ID = types.String{Value: strconv.Itoa(newIssueScreen.ID)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r jiraIssueScreenResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Reading issue screen resource")

	var state jiraIssueScreenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue screen from state", map[string]interface{}{
		"issueScreenState": fmt.Sprintf("%+v", state),
	})

	issueScreenId, _ := strconv.Atoi(state.ID.Value)

	resIssueScreen, res, err := r.p.jira.Screen.Gets(ctx, []int{issueScreenId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue screen, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue screen from API state")

	state.Name = types.String{Value: resIssueScreen.Values[0].Name}
	state.Description = types.String{Value: resIssueScreen.Values[0].Description}
	tflog.Debug(ctx, "Updated state with API state")

	tflog.Debug(ctx, "Storing issue screen info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r jiraIssueScreenResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating issue screen")

	var plan jiraIssueScreenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded issue screen configuration", map[string]interface{}{
		"issueScreenConfig": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueScreenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded issue screen from state", map[string]interface{}{
		"issueScreenState": fmt.Sprintf("%+v", state),
	})

	issueScreenId, _ := strconv.Atoi(state.ID.Value)
	_, res, err := r.p.jira.Screen.Update(ctx, issueScreenId, plan.Name.Value, plan.Description.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue screen, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Updated issue screen in API state")

	var updatedState = jiraIssueScreenResourceModel{
		ID:          types.String{Value: state.ID.Value},
		Name:        types.String{Value: plan.Name.Value},
		Description: types.String{Value: plan.Description.Value},
	}

	tflog.Debug(ctx, "Storing issue screen info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r jiraIssueScreenResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Deleting issue screen resource")

	var state jiraIssueScreenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue screen from state")

	issueScreenId, _ := strconv.Atoi(state.ID.Value)
	res, err := r.p.jira.Screen.Delete(ctx, issueScreenId)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue screen, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Removed issue screen from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
