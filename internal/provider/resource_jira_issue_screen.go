package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/planmodifiers/stringmodifiers"
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

func (*jiraIssueScreenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Screen Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the issue screen.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the screen." +
					"The name must be unique." +
					"The maximum length is 255 characters.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the screen." +
					"The maximum length is 255 characters.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				PlanModifiers: []planmodifier.String{
					stringmodifiers.DefaultValue(""),
				},
			},
		},
	}
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
	tflog.Debug(ctx, "Loaded issue screen plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	newIssueScreen, res, err := r.p.jira.Screen.Create(ctx, plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue screen, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue screen")

	plan.ID = types.StringValue(strconv.Itoa(newIssueScreen.ID))

	tflog.Debug(ctx, "Storing issue screen info into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
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
		"readState": fmt.Sprintf("%+v", state),
	})

	issueScreenId, _ := strconv.Atoi(state.ID.ValueString())

	issueScreen, res, err := r.p.jira.Screen.Gets(ctx, []int{issueScreenId}, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue screen, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue screen from API state")

	state.Name = types.StringValue(issueScreen.Values[0].Name)
	state.Description = types.StringValue(issueScreen.Values[0].Description)

	tflog.Debug(ctx, "Storing issue screen info into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueScreenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue screen resource")

	var plan jiraIssueScreenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue screen plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueScreenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue screen from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

	issueScreenId, _ := strconv.Atoi(state.ID.ValueString())
	_, res, err := r.p.jira.Screen.Update(ctx, issueScreenId, plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue screen, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Updated issue screen in API state")

	var updatedState = jiraIssueScreenResourceModel{
		ID:          types.StringValue(state.ID.ValueString()),
		Name:        types.StringValue(plan.Name.ValueString()),
		Description: types.StringValue(plan.Description.ValueString()),
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

	issueScreenId, _ := strconv.Atoi(state.ID.ValueString())
	res, err := r.p.jira.Screen.Delete(ctx, issueScreenId)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue screen, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Removed issue screen from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
