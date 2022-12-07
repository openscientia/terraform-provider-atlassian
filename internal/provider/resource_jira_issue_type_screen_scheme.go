package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
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

func (*jiraIssueTypeScreenSchemeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Type Screen Scheme Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the issue type screen scheme.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the issue type screen scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the issue type screen scheme. " +
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
			"issue_type_mappings": schema.ListNestedAttribute{
				MarkdownDescription: "The IDs of the screen schemes for the issue type IDs and default. " +
					"A default entry is required to create an issue type screen scheme, it defines the mapping for all issue types without a screen scheme.",
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"issue_type_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the issue type or default. " +
								"Only issue types used in classic projects are accepted. " +
								"An entry for default must be provided and defines the mapping for all issue types without a screen scheme.",
							Required: true,
						},
						"screen_scheme_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the screen scheme. " +
								"Only screen schemes used in classic projects are accepted.",
							Required: true,
						},
					},
				},
			},
		},
	}
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
	tflog.Debug(ctx, "Creating issue type screen scheme resource")

	var plan jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	issueTypeMappings := []*models.IssueTypeScreenSchemeMappingPayloadScheme{}
	for _, v := range plan.IssueTypeMappings {
		issueTypeMappings = append(issueTypeMappings, &models.IssueTypeScreenSchemeMappingPayloadScheme{
			IssueTypeID:    v.IssueTypeId.ValueString(),
			ScreenSchemeID: v.ScreenSchemeId.ValueString(),
		})
	}

	createRequestPayload := models.IssueTypeScreenSchemePayloadScheme{
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		IssueTypeMappings: issueTypeMappings,
	}

	newIssueTypeScreenScheme, res, err := r.p.jira.Issue.Type.ScreenScheme.Create(ctx, &createRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue type screen scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue type screen scheme")

	plan.ID = types.StringValue(newIssueTypeScreenScheme.ID)

	tflog.Debug(ctx, "Storing issue type screen scheme into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueTypeScreenSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type screen scheme resource")

	var state jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	issueTypeScreenSchemeId, _ := strconv.Atoi(state.ID.ValueString())
	options := &models.ScreenSchemeParamsScheme{
		IDs: []int{issueTypeScreenSchemeId},
	}
	issueTypeScreenSchemeDetails, res, err := r.p.jira.Issue.Type.ScreenScheme.Gets(ctx, options, 0, 1)
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

	state.Name = types.StringValue(issueTypeScreenSchemeDetails.Values[0].Name)
	state.Description = types.StringValue(issueTypeScreenSchemeDetails.Values[0].Description)
	var mappings []jiraIssueTypeScreenSchemeMapping
	for _, v := range issueTypeScreenSchemeMappings.Values {
		mappings = append(mappings, jiraIssueTypeScreenSchemeMapping{
			IssueTypeId:    types.StringValue(v.IssueTypeID),
			ScreenSchemeId: types.StringValue(v.ScreenSchemeID),
		})
	}
	state.IssueTypeMappings = mappings

	tflog.Debug(ctx, "Storing issue type screen scheme into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueTypeScreenSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating issue type screen scheme resource")

	var plan jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme plan", map[string]interface{}{
		"updatePlan": fmt.Sprintf("%+v", plan),
	})

	var state jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme from state", map[string]interface{}{
		"updateState": fmt.Sprintf("%+v", state),
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

	plan.ID = types.StringValue(state.ID.ValueString())

	tflog.Debug(ctx, "Storing issue type screen scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueTypeScreenSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting issue type screen scheme resource")

	var state jiraIssueTypeScreenSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme from state")

	res, err := r.p.jira.Issue.Type.ScreenScheme.Delete(ctx, state.ID.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete issue type screen scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted issue type screen scheme from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}

func (r *jiraIssueTypeScreenSchemeResource) updateNameAndDescription(ctx context.Context, p, s *jiraIssueTypeScreenSchemeResourceModel) error {
	if p.Name.ValueString() != s.Name.ValueString() || p.Description.ValueString() != s.Description.ValueString() {
		res, err := r.p.jira.Issue.Type.ScreenScheme.Update(ctx, s.ID.ValueString(), p.Name.ValueString(), p.Description.ValueString())
		if err != nil {
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			return fmt.Errorf(" Unable to update issue type screen scheme name and description, got error: %s\n%s", err, resBody)
		}
		tflog.Debug(ctx, "Updated issue type screen scheme name and description", map[string]interface{}{
			"newNameAndDescription": fmt.Sprintf("%s, %s", p.Name.ValueString(), p.Description.ValueString()),
		})
	}

	return nil
}

func (r *jiraIssueTypeScreenSchemeResource) updateDefaultMapping(ctx context.Context, p, s *jiraIssueTypeScreenSchemeResourceModel) error {
	var planDefaultMapping *jiraIssueTypeScreenSchemeMapping
	for _, m := range p.IssueTypeMappings {
		if m.IssueTypeId.ValueString() == "default" {
			planDefaultMapping = &jiraIssueTypeScreenSchemeMapping{
				IssueTypeId:    m.IssueTypeId,
				ScreenSchemeId: m.ScreenSchemeId,
			}
		}
	}
	for _, m := range s.IssueTypeMappings {
		if m.IssueTypeId.ValueString() == "default" && m.ScreenSchemeId.ValueString() != planDefaultMapping.ScreenSchemeId.ValueString() {
			res, err := r.p.jira.Issue.Type.ScreenScheme.UpdateDefault(ctx, s.ID.ValueString(), planDefaultMapping.ScreenSchemeId.ValueString())
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
			if pm.IssueTypeId.ValueString() == "default" || pm == sm {
				canAdd = false
			}
		}
		if canAdd {
			addMappingPayload := &models.IssueTypeScreenSchemePayloadScheme{
				IssueTypeMappings: []*models.IssueTypeScreenSchemeMappingPayloadScheme{
					{
						IssueTypeID:    pm.IssueTypeId.ValueString(),
						ScreenSchemeID: pm.ScreenSchemeId.ValueString(),
					},
				},
			}
			res, err := r.p.jira.Issue.Type.ScreenScheme.Append(ctx, s.ID.ValueString(), addMappingPayload)
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
			if sm.IssueTypeId.ValueString() == "default" || sm == pm {
				canRemove = false
			}
		}
		if canRemove {
			removeMappings = append(removeMappings, sm.IssueTypeId.ValueString())
		}
	}
	if len(removeMappings) > 0 {
		res, err := r.p.jira.Issue.Type.ScreenScheme.Remove(ctx, s.ID.ValueString(), removeMappings)
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
