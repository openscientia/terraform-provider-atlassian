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
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_plan_modification"
)

type (
	jiraIssueTypeScreenSchemeResource struct {
		p atlassianProvider
	}

	jiraIssueTypeScreenSchemeResourceModel struct {
		ID                types.String                       `tfsdk:"id"`
		Name              types.String                       `tfsdk:"name"`
		Description       types.String                       `tfsdk:"description"`
		IssueTypeMappings []jiraIssueTypeScreenSchemeMapping `tfsdk:"issue_type_mappings"`
	}

	jiraIssueTypeScreenSchemeMapping struct {
		IssueTypeId    types.String `tfsdk:"issue_type_id"`
		ScreenSchemeId types.String `tfsdk:"screen_scheme_id"`
	}
)

var (
	_ resource.Resource                = (*jiraIssueTypeScreenSchemeResource)(nil)
	_ resource.ResourceWithImportState = (*jiraIssueTypeScreenSchemeResource)(nil)
)

func NewJiraIssueTypeScreenSchemeResource() resource.Resource {
	return &jiraIssueTypeScreenSchemeResource{}
}

func (*jiraIssueTypeScreenSchemeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_type_screen_scheme"
}

func (*jiraIssueTypeScreenSchemeResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Type Screen Scheme Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type screen scheme.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue type screen scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the issue type screen scheme. " +
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
			"issue_type_mappings": {
				MarkdownDescription: "The IDs of the screen schemes for the issue type IDs and default. " +
					"A default entry is required to create an issue type screen scheme, it defines the mapping for all issue types without a screen scheme.",
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"issue_type_id": {
							MarkdownDescription: "The ID of the issue type or default. " +
								"Only issue types used in classic projects are accepted. " +
								"An entry for default must be provided and defines the mapping for all issue types without a screen scheme.",
							Required: true,
							Type:     types.StringType,
						},
						"screen_scheme_id": {
							MarkdownDescription: "The ID of the screen scheme. " +
								"Only screen schemes used in classic projects are accepted.",
							Required: true,
							Type:     types.StringType,
						},
					},
				),
			},
		},
	}, nil
}

func (r *jiraIssueTypeScreenSchemeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraIssueTypeScreenSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraIssueTypeScreenSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating issue type screen scheme")

	var plan jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme configuration", map[string]interface{}{
		"issueTypeScreenSchemeConfig": fmt.Sprintf("%+v", plan),
	})

	issueTypeMappings := []*models.IssueTypeScreenSchemeMappingPayloadScheme{}
	for _, v := range plan.IssueTypeMappings {
		issueTypeMappings = append(issueTypeMappings, &models.IssueTypeScreenSchemeMappingPayloadScheme{
			IssueTypeID:    v.IssueTypeId.Value,
			ScreenSchemeID: v.ScreenSchemeId.Value,
		})
	}

	createRequestPayload := models.IssueTypeScreenSchemePayloadScheme{
		Name:              plan.Name.Value,
		IssueTypeMappings: issueTypeMappings,
	}
	tflog.Debug(ctx, "Generated request payload", map[string]interface{}{
		"issueTypeScreenScheme": fmt.Sprintf("%+v", createRequestPayload),
	})

	newIssueTypeScreenScheme, res, err := r.p.jira.Issue.Type.ScreenScheme.Create(ctx, &createRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue type screen scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue type screen scheme", map[string]interface{}{
		"issueTypeScreenScheme": fmt.Sprintf("%+v", newIssueTypeScreenScheme),
	})

	plan.ID = types.String{Value: newIssueTypeScreenScheme.ID}

	// TODO: Remove this when 'description' can be addded on create call above
	// https://github.com/ctreminiom/go-atlassian/issues/131
	res, err = r.p.jira.Issue.Type.ScreenScheme.Update(ctx, newIssueTypeScreenScheme.ID, plan.Name.Value, plan.Description.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add issue type screen scheme description, got error: %s\n%s", err, resBody))
		return
	}

	tflog.Debug(ctx, "Storing issue type screen scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueTypeScreenSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type screen scheme")

	var state jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme from state", map[string]interface{}{
		"issueTypeScreenSchemeState": fmt.Sprintf("%+v", state),
	})

	issueTypeScreenSchemeId, _ := strconv.Atoi(state.ID.Value)
	issueTypeScreenSchemeDetails, res, err := r.p.jira.Issue.Type.ScreenScheme.Gets(ctx, []int{issueTypeScreenSchemeId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type screen scheme, got error: %s\n%s", err, resBody))
		return
	}

	issueTypeScreenSchemeMappings, res, err := r.p.jira.Issue.Type.ScreenScheme.Mapping(ctx, []int{issueTypeScreenSchemeId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type screen scheme mappings, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue type screen scheme from API state")

	state.Name = types.String{Value: issueTypeScreenSchemeDetails.Values[0].Name}
	state.Description = types.String{Value: issueTypeScreenSchemeDetails.Values[0].Description}
	var mappings []jiraIssueTypeScreenSchemeMapping
	for _, v := range issueTypeScreenSchemeMappings.Values {
		mappings = append(mappings, jiraIssueTypeScreenSchemeMapping{
			IssueTypeId:    types.String{Value: v.IssueTypeID},
			ScreenSchemeId: types.String{Value: v.ScreenSchemeID},
		})
	}
	state.IssueTypeMappings = mappings

	tflog.Debug(ctx, "Storing issue type screen scheme into the state", map[string]interface{}{
		"newState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueTypeScreenSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue type screen scheme")

	var plan jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme configuration", map[string]interface{}{
		"issueTypeScreenSchemeConfig": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme from state", map[string]interface{}{
		"issueTypeScreenSchemeState": fmt.Sprintf("%+v", state),
	})

	err := r.updateNameAndDescription(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	err = r.updateDefaultMapping(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	err = r.addMappings(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	err = r.removeMappings(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Updated issue type screen scheme in API state")

	plan.ID = types.String{Value: state.ID.Value}

	tflog.Debug(ctx, "Storing issue type screen scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueTypeScreenSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue type screen scheme")

	var state jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme from state")

	res, err := r.p.jira.Issue.Type.ScreenScheme.Delete(ctx, state.ID.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue type screen scheme, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Removed issue type screen scheme from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}

func (r *jiraIssueTypeScreenSchemeResource) updateNameAndDescription(ctx context.Context, p, s *jiraIssueTypeScreenSchemeResourceModel) error {
	if p.Name.Value != s.Name.Value || p.Description.Value != s.Description.Value {
		res, err := r.p.jira.Issue.Type.ScreenScheme.Update(ctx, s.ID.Value, p.Name.Value, p.Description.Value)
		if err != nil {
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			return fmt.Errorf(" Unable to update issue type screen scheme name and description, got error: %s\n%s", err, resBody)
		}
		tflog.Debug(ctx, "Updated issue type screen scheme name and description", map[string]interface{}{
			"newNameAndDescription": fmt.Sprintf("%s, %s", p.Name.Value, p.Description.Value),
		})
	}

	return nil
}

func (r *jiraIssueTypeScreenSchemeResource) updateDefaultMapping(ctx context.Context, p, s *jiraIssueTypeScreenSchemeResourceModel) error {
	var planDefaultMapping *jiraIssueTypeScreenSchemeMapping
	for _, m := range p.IssueTypeMappings {
		if m.IssueTypeId.Value == "default" {
			planDefaultMapping = &jiraIssueTypeScreenSchemeMapping{
				IssueTypeId:    m.IssueTypeId,
				ScreenSchemeId: m.ScreenSchemeId,
			}
		}
	}
	for _, m := range s.IssueTypeMappings {
		if m.IssueTypeId.Value == "default" && m.ScreenSchemeId.Value != planDefaultMapping.ScreenSchemeId.Value {
			res, err := r.p.jira.Issue.Type.ScreenScheme.UpdateDefault(ctx, s.ID.Value, planDefaultMapping.ScreenSchemeId.Value)
			if err != nil {
				var resBody string
				if res != nil {
					return fmt.Errorf(" Unable to update issue type screen scheme default mapping, got error: %s\n%s", err, resBody)
				}
			}
			tflog.Debug(ctx, "Updated issue type screen scheme default mapping", map[string]interface{}{
				"newDefaultMapping": fmt.Sprintf("%+v", planDefaultMapping),
			})
		}
	}

	return nil
}

func (r *jiraIssueTypeScreenSchemeResource) addMappings(ctx context.Context, p, s *jiraIssueTypeScreenSchemeResourceModel) error {
	var canAdd bool
	for _, pm := range p.IssueTypeMappings {
		canAdd = true
		for _, sm := range s.IssueTypeMappings {
			// Skip default mapping or existing mapping in state
			if pm.IssueTypeId.Value == "default" || pm == sm {
				canAdd = false
			}
		}
		if canAdd {
			addMappingPayload := &models.IssueTypeScreenSchemePayloadScheme{
				IssueTypeMappings: []*models.IssueTypeScreenSchemeMappingPayloadScheme{
					{
						IssueTypeID:    pm.IssueTypeId.Value,
						ScreenSchemeID: pm.ScreenSchemeId.Value,
					},
				},
			}
			res, err := r.p.jira.Issue.Type.ScreenScheme.Append(ctx, s.ID.Value, addMappingPayload)
			if err != nil {
				var resBody string
				if res != nil {
					resBody = res.Bytes.String()
				}
				return fmt.Errorf(" Unable to add issue type screen scheme mapping, got error: %s\n%s", err, resBody)
			}
			tflog.Debug(ctx, "Added issue type screen scheme mapping", map[string]interface{}{
				"newMapping": fmt.Sprintf("%+v", *addMappingPayload.IssueTypeMappings[0]),
			})
		}
	}

	return nil
}

func (r *jiraIssueTypeScreenSchemeResource) removeMappings(ctx context.Context, p, s *jiraIssueTypeScreenSchemeResourceModel) error {
	var removeMappings []string
	var canRemove bool
	for _, sm := range s.IssueTypeMappings {
		canRemove = true
		for _, pm := range p.IssueTypeMappings {
			// Skip default mapping or existing mapping in plan
			if sm.IssueTypeId.Value == "default" || sm == pm {
				canRemove = false
			}
		}
		if canRemove {
			removeMappings = append(removeMappings, sm.IssueTypeId.Value)
		}
	}
	if len(removeMappings) > 0 {
		res, err := r.p.jira.Issue.Type.ScreenScheme.Remove(ctx, s.ID.Value, removeMappings)
		if err != nil {
			var resBody string
			if res != nil {
				return fmt.Errorf(" Unable to remove issue type screen scheme mappings, got error: %s\n%s", err, resBody)
			}
		}
		tflog.Debug(ctx, "Removed issue type screen scheme mappings", map[string]interface{}{
			"removedMappings": fmt.Sprintf("%+v", removeMappings),
		})
	}

	return nil
}
