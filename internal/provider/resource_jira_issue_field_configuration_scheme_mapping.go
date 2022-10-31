package atlassian

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueFieldConfigurationSchemeMappingResource struct {
		p atlassianProvider
	}

	jiraIssueFieldConfigurationSchemeMappingResourceModel struct {
		ID                         types.String `tfsdk:"id"`
		FieldConfigurationSchemeID types.String `tfsdk:"field_configuration_scheme_id"`
		FieldConfigurationID       types.String `tfsdk:"field_configuration_id"`
		IssueTypeID                types.String `tfsdk:"issue_type_id"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueFieldConfigurationSchemeMappingResource)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueFieldConfigurationSchemeMappingResource)(nil)
)

func NewJiraIssueFieldConfigurationSchemeMappingResource() resource.Resource {
	return &jiraIssueFieldConfigurationSchemeMappingResource{}
}

func (*jiraIssueFieldConfigurationSchemeMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_field_configuration_scheme_mapping"
}

func (*jiraIssueFieldConfigurationSchemeMappingResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Field Configuration Scheme Mapping Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue field configuration scheme mapping. " +
					"It is computed using `field_configuration_scheme_id`, `field_configuration_id` and `issue_type_id` separated by a hyphen (`-`).",
				Computed: true,
				Type:     types.StringType,
			},
			"field_configuration_scheme_id": {
				MarkdownDescription: "(Forces new resource) The ID of the issue field configuration scheme.",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"field_configuration_id": {
				MarkdownDescription: "(Forces new resource) The ID of the issue field configuration.",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"issue_type_id": {
				MarkdownDescription: "(Forces new resource) The ID of the issue type or `default`. " +
					"When set to `default` this issue field configuration scheme mapping applies to all issue types without an issue field configuration.",
				Required: true,
				Type:     types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *jiraIssueFieldConfigurationSchemeMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraIssueFieldConfigurationSchemeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError("Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: field_configuration_scheme_id, field_configuration_id, issue_type_id. Got: %q", req.ID))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Importing issue field configuration scheme mapping with import identifier: %+v", idParts))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("field_configuration_scheme_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("field_configuration_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("issue_type_id"), idParts[2])...)
}

func (r *jiraIssueFieldConfigurationSchemeMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue field configuration scheme mapping resource")

	var plan jiraIssueFieldConfigurationSchemeMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme mapping plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	issueFieldConfigurationSchemeId, _ := strconv.Atoi(plan.FieldConfigurationSchemeID.ValueString())
	createRequestPayload := models.FieldConfigurationToIssueTypeMappingPayloadScheme{
		Mappings: []*models.FieldConfigurationToIssueTypeMappingScheme{
			{
				IssueTypeID:          plan.IssueTypeID.ValueString(),
				FieldConfigurationID: plan.FieldConfigurationID.ValueString(),
			},
		},
	}

	res, err := r.p.jira.Issue.Field.Configuration.Scheme.Link(ctx, issueFieldConfigurationSchemeId, &createRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue field configuration scheme mapping, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue field configuration scheme mapping")

	plan.ID = types.String{Value: createIssueFieldConfigurationSchemeMappingID(plan.FieldConfigurationSchemeID.ValueString(), plan.FieldConfigurationID.ValueString(), plan.IssueTypeID.ValueString())}

	tflog.Debug(ctx, "Storing issue field configuration scheme mapping into the state", map[string]interface{}{
		"newState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationSchemeMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue field configuration scheme mapping resource")

	var state jiraIssueFieldConfigurationSchemeMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme mapping from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	fieldConfigurationSchemeId, _ := strconv.Atoi(state.FieldConfigurationSchemeID.ValueString())
	mappings, res, err := r.p.jira.Issue.Field.Configuration.Scheme.Mapping(ctx, []int{fieldConfigurationSchemeId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration scheme mappings, got error: %s\n%s", err, resBody))
		return
	}

	found := false
	for _, m := range mappings.Values {
		if m.FieldConfigurationSchemeID == state.FieldConfigurationSchemeID.ValueString() {
			if m.IssueTypeID == state.IssueTypeID.ValueString() && m.FieldConfigurationID == state.FieldConfigurationID.ValueString() {
				found = true
			}
		}
	}

	if !found {
		// If mapping not found in API state it means that resource was changed outside Terraform
		// and it must be deleted from Terraform state and recreated
		tflog.Warn(ctx, "Unable to find issue field configuration scheme mapping in API state, deleting resource from state")
		resp.State.RemoveResource(ctx)
		return
	}
	tflog.Debug(ctx, "Retrieved issue field configuration scheme mapping from API state")

	state.ID = types.String{Value: createIssueFieldConfigurationSchemeMappingID(state.FieldConfigurationSchemeID.ValueString(), state.FieldConfigurationID.ValueString(), state.IssueTypeID.ValueString())}

	tflog.Debug(ctx, "Storing issue field configuration scheme mapping into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueFieldConfigurationSchemeMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The RequiresReplace plan modifier will trigger Terraform to destroy and recreate the resource
	// if any of the required attributes changes, i.e. field_configuration_scheme_id, field_scheme_id or issue_type_id
	tflog.Debug(ctx, "If the value of any required attribute changes, Terraform will destroy and recreate the resource")
}

func (r *jiraIssueFieldConfigurationSchemeMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue field configuration scheme mapping resource")

	var state jiraIssueFieldConfigurationSchemeMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme mapping from state")

	if state.IssueTypeID.ValueString() == "default" {
		// It is not possible to delete a "default" mapping from API state
		// because field configuration schemes must have at least a "default mapping
		// therefore, the resource will only be deleted from Terraform state
		return
	}

	fieldConfigurationSchemeId, _ := strconv.Atoi(state.FieldConfigurationSchemeID.ValueString())
	res, err := r.p.jira.Issue.Field.Configuration.Scheme.Unlink(ctx, fieldConfigurationSchemeId, []string{state.IssueTypeID.ValueString()})
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue field configuration scheme mapping, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted issue field configuration scheme mapping from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}

func createIssueFieldConfigurationSchemeMappingID(fieldConfigurationSchemeId, fieldConfigurationId, issueTypeId string) string {
	return strings.Join([]string{fieldConfigurationSchemeId, fieldConfigurationId, issueTypeId}, "-")
}
