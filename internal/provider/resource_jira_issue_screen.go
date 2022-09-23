package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_validation"
)

type (
	jiraIssueScreenResource struct {
		p atlassianProvider
	}

	jiraIssueScreenResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueScreenResource)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueScreenResource)(nil)
)

func NewJiraIssueScreenResource() resource.Resource {
	return &jiraIssueScreenResource{}
}

func (*jiraIssueScreenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_screen"
}

func (*jiraIssueScreenResource) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r *jiraIssueScreenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraIssueScreenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraIssueScreenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue screen resource")

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

func (r *jiraIssueScreenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (r *jiraIssueScreenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r *jiraIssueScreenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
