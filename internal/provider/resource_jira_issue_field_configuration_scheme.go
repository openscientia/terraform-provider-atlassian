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
	jiraIssueFieldConfigurationSchemeResource struct {
		p atlassianProvider
	}

	jiraIssueFieldConfigurationSchemeResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueFieldConfigurationSchemeResource)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueFieldConfigurationSchemeResource)(nil)
)

func NewJiraIssueFieldConfigurationSchemeResource() resource.Resource {
	return &jiraIssueFieldConfigurationSchemeResource{}
}

func (*jiraIssueFieldConfigurationSchemeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_field_configuration_scheme"
}

func (*jiraIssueFieldConfigurationSchemeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Field Configuration Scheme Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the issue field configuration scheme.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the issue field configuration scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the issue field configuration scheme. " +
					"The maximum length is 1024 characters.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
				PlanModifiers: []planmodifier.String{
					stringmodifiers.DefaultValue(""),
				},
			},
		},
	}
}

func (r *jiraIssueFieldConfigurationSchemeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraIssueFieldConfigurationSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue field configuration scheme resource")

	var plan jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	issueFieldConfigurationScheme, res, err := r.p.jira.Issue.Field.Configuration.Scheme.Create(ctx, plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue field configuration scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue field configuration scheme")

	plan.ID = types.StringValue(issueFieldConfigurationScheme.ID)

	tflog.Debug(ctx, "Storing issue field configuration scheme info into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue field configuration scheme resource")

	var state jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	id, _ := strconv.Atoi(state.ID.ValueString())
	issueFieldConfigurationScheme, res, err := r.p.jira.Issue.Field.Configuration.Scheme.Gets(ctx, []int{id}, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue field configuration scheme from API state")

	state.Name = types.StringValue(issueFieldConfigurationScheme.Values[0].Name)
	state.Description = types.StringValue(issueFieldConfigurationScheme.Values[0].Description)

	tflog.Debug(ctx, "Storing issue field configuration scheme info into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue field configuration scheme resource")

	var plan jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
	})

	id, _ := strconv.Atoi(state.ID.ValueString())
	res, err := r.p.jira.Issue.Field.Configuration.Scheme.Update(ctx, id, plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue field configuration scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Updated issue field configuration scheme")

	plan.ID = types.StringValue(state.ID.ValueString())

	tflog.Debug(ctx, "Storing issue field configuration scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue field configuration scheme resource")

	var state jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme from state")

	id, _ := strconv.Atoi(state.ID.ValueString())
	res, err := r.p.jira.Issue.Field.Configuration.Scheme.Delete(ctx, id)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue field configuration scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted issue field configuration scheme from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
