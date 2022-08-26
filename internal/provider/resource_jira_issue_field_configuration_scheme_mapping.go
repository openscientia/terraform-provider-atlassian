package atlassian

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueFieldConfigurationSchemeMappingResource struct {
		p atlassianProvider
	}
	jiraIssueFieldConfigurationSchemeMappingResourceType struct{}

	jiraIssueFieldConfigurationSchemeMappingResourceModel struct {
		ID                         types.String `tfsdk:"id"`
		FieldConfigurationSchemeID types.String `tfsdk:"field_configuration_scheme_id"`
		FieldConfigurationID       types.String `tfsdk:"field_configuration_id"`
		IssueTypeID                types.String `tfsdk:"issue_type_id"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueFieldConfigurationSchemeMappingResource)(nil)
	_ provider.ResourceType            = (*jiraIssueFieldConfigurationSchemeMappingResourceType)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueFieldConfigurationSchemeMappingResource)(nil)
)

func (*jiraIssueFieldConfigurationSchemeMappingResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r *jiraIssueFieldConfigurationSchemeMappingResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraIssueFieldConfigurationSchemeMappingResource{
		p: provider,
	}, diags
}

func (r *jiraIssueFieldConfigurationSchemeMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
	}

	var plan jiraIssueFieldConfigurationSchemeMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme mapping plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	issueFieldConfigurationSchemeId, _ := strconv.Atoi(plan.FieldConfigurationSchemeID.Value)
	createRequestPayload := models.FieldConfigurationToIssueTypeMappingPayloadScheme{
		Mappings: []*models.FieldConfigurationToIssueTypeMappingScheme{
			{
				IssueTypeID:          plan.IssueTypeID.Value,
				FieldConfigurationID: plan.FieldConfigurationID.Value,
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

	plan.ID = types.String{Value: createIssueFieldConfigurationSchemeMappingID(plan.FieldConfigurationSchemeID.Value, plan.FieldConfigurationID.Value, plan.IssueTypeID.Value)}

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

	fieldConfigurationSchemeId, _ := strconv.Atoi(state.FieldConfigurationSchemeID.Value)
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
		if m.FieldConfigurationSchemeID == state.FieldConfigurationSchemeID.Value {
			if m.IssueTypeID == state.IssueTypeID.Value && m.FieldConfigurationID == state.FieldConfigurationID.Value {
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

	state.ID = types.String{Value: createIssueFieldConfigurationSchemeMappingID(state.FieldConfigurationSchemeID.Value, state.FieldConfigurationID.Value, state.IssueTypeID.Value)}

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

	if state.IssueTypeID.Value == "default" {
		// It is not possible to delete a "default" mapping from API state
		// because field configuration schemes must have at least a "default mapping
		// therefore, the resource will only be deleted from Terraform state
		return
	}

	fieldConfigurationSchemeId, _ := strconv.Atoi(state.FieldConfigurationSchemeID.Value)
	res, err := r.p.jira.Issue.Field.Configuration.Scheme.Unlink(ctx, fieldConfigurationSchemeId, []string{state.IssueTypeID.Value})
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
