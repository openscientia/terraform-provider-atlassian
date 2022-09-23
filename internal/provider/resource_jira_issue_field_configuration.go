package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
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
	jiraIssueFieldConfigurationResource struct {
		p atlassianProvider
	}

	jiraIssueFieldConfigurationResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueFieldConfigurationResource)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueFieldConfigurationResource)(nil)
)

func NewJiraIssueFieldConfigurationResource() resource.Resource {
	return &jiraIssueFieldConfigurationResource{}
}

func (*jiraIssueFieldConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_field_configuration"
}

func (*jiraIssueFieldConfigurationResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Field Configuration Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue field configuration.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue field configuration. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the issue field configuration. " +
					"The maximum length is 255 characters.",
				Optional: true,
				Computed: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(255),
				},
				PlanModifiers: tfsdk.AttributePlanModifiers{
					attribute_plan_modification.DefaultValue(types.String{Value: ""}),
				},
			},
		},
	}, nil
}

func (r *jiraIssueFieldConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraIssueFieldConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraIssueFieldConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue field configuration")

	var plan jiraIssueFieldConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration plan", map[string]interface{}{
		"issueFieldConfigurationPlan": fmt.Sprintf("%+v", plan),
	})

	issueFieldConfiguration, res, err := r.p.jira.Issue.Field.Configuration.Create(ctx, plan.Name.Value, plan.Description.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue field configuration, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue field configuration", map[string]interface{}{
		"issueFieldConfiguration": fmt.Sprintf("%+v", issueFieldConfiguration),
	})

	plan.ID = types.String{Value: strconv.Itoa(issueFieldConfiguration.ID)}

	tflog.Debug(ctx, "Storing issue field configuration info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue field configuration")

	var state jiraIssueFieldConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration from state", map[string]interface{}{
		"issueFieldConfigurationState": fmt.Sprintf("%+v", state),
	})

	issueFieldConfigurationId, _ := strconv.Atoi(state.ID.Value)
	issueFieldConfiguration, res, err := r.p.jira.Issue.Field.Configuration.Gets(ctx, []int{issueFieldConfigurationId}, false, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue field configuration from API state")

	state.Name = types.String{Value: issueFieldConfiguration.Values[0].Name}
	state.Description = types.String{Value: issueFieldConfiguration.Values[0].Description}

	tflog.Debug(ctx, "Storing issue field configuration into the state", map[string]interface{}{
		"issueFieldConfiguration": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueFieldConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue field configuration")

	var plan jiraIssueFieldConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration plan", map[string]interface{}{
		"issueFieldConfigurationPlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueFieldConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration from state", map[string]interface{}{
		"issueFieldConfigurationState": fmt.Sprintf("%+v", state),
	})

	issueFieldConfigurationId, _ := strconv.Atoi(state.ID.Value)
	res, err := r.p.jira.Issue.Field.Configuration.Update(ctx, issueFieldConfigurationId, plan.Name.Value, plan.Description.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue field configuration, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Updated issue field configuration in API state")

	plan.ID = types.String{Value: state.ID.Value}

	tflog.Debug(ctx, "Storing issue field configuration info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue field configuration")

	var state jiraIssueFieldConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration from state")

	issueFieldConfigurationID, _ := strconv.Atoi(state.ID.Value)
	res, err := r.p.jira.Issue.Field.Configuration.Delete(ctx, issueFieldConfigurationID)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue field configuration, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Removed issue field configuration from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
