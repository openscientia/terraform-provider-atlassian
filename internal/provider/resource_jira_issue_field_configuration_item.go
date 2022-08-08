package atlassian

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueFieldConfigurationItemResource struct {
		p provider
	}

	jiraIssueFieldConfigurationItemResourceType struct{}

	jiraIssueFieldConfigurationItemResourceModel struct {
		ID                      types.String                     `tfsdk:"id"`
		IssueFieldConfiguration types.String                     `tfsdk:"issue_field_configuration"`
		Item                    *jiraIssueFieldConfigurationItem `tfsdk:"item"`
	}

	jiraIssueFieldConfigurationItem struct {
		ID          types.String `tfsdk:"id"`
		Description types.String `tfsdk:"description"`
		IsHidden    types.Bool   `tfsdk:"is_hidden"`
		IsRequired  types.Bool   `tfsdk:"is_required"`
		Renderer    types.String `tfsdk:"renderer"`
	}
)

var (
	_                   tfsdk.Resource                = (*jiraIssueFieldConfigurationItemResource)(nil)
	_                   tfsdk.ResourceType            = (*jiraIssueFieldConfigurationItemResourceType)(nil)
	_                   tfsdk.ResourceWithImportState = (*jiraIssueFieldConfigurationItemResource)(nil)
	renderableItemTypes                               = []string{"string", "comments-page"}
)

func (*jiraIssueFieldConfigurationItemResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Field Configuration Item Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue field configuration item. " +
					"It is computed using `issue_field_configuration` and `item.id` separated by a hyphen (`-`).",
				Computed: true,
				Type:     types.StringType,
			},
			"issue_field_configuration": {
				MarkdownDescription: "(Forces new resource) The ID of the issue field configuration.",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"item": {
				MarkdownDescription: "Details of a field within the issue field configuration.",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"id": {
							MarkdownDescription: "(Forces new resource) The ID of the field within the issue field configuration.",
							Required:            true,
							Type:                types.StringType,
							Validators: []tfsdk.AttributeValidator{
								stringvalidator.RegexMatches(regexp.MustCompile(`^customfield_[0-9]{5}$|^[a-zA-Z]*$`), ""),
							},
							PlanModifiers: tfsdk.AttributePlanModifiers{
								tfsdk.RequiresReplace(),
							},
						},
						"description": {
							MarkdownDescription: "The description of the field within the issue field configuration.",
							Computed:            true,
							Optional:            true,
							Type:                types.StringType,
						},
						"is_hidden": {
							MarkdownDescription: "Whether the field is hidden in the issue field configuration. " +
								"Can be `true` or `false`.",
							Computed: true,
							Optional: true,
							Type:     types.BoolType,
						},
						"is_required": {
							MarkdownDescription: "Whether the field is required in the issue field configuration. " +
								"Can be `true` or `false`.",
							Computed: true,
							Optional: true,
							Type:     types.BoolType,
						},
						"renderer": {
							MarkdownDescription: "The renderer type for the field within the issue field configuration. " +
								"Can be `text-renderer` or `wiki-renderer`.",
							Computed: true,
							Optional: true,
							Type:     types.StringType,
							Validators: []tfsdk.AttributeValidator{
								stringvalidator.OneOf("text-renderer", "wiki-renderer"),
							},
						},
					},
				),
			},
		},
	}, nil
}

func (r *jiraIssueFieldConfigurationItemResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraIssueFieldConfigurationItemResource{
		p: provider,
	}, diags
}

func (r *jiraIssueFieldConfigurationItemResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: issue_field_configuration, item.id. Got: %q", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("issue_field_configuration"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("item").AtName("id"), idParts[1])...)
}

func (r *jiraIssueFieldConfigurationItemResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating issue field configuration item")

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
	}

	var plan jiraIssueFieldConfigurationItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration item plan", map[string]interface{}{
		"issueFieldConfigurationItem": fmt.Sprintf("%+v, %+v", plan, *plan.Item),
	})

	if !plan.Item.Renderer.IsNull() && !plan.Item.Renderer.IsUnknown() {
		err := r.checkIssueFieldConfigurationItemRenderable(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
	}

	issueFieldConfigurationId, _ := strconv.Atoi(plan.IssueFieldConfiguration.Value)
	createRequestPayload := models.UpdateFieldConfigurationItemPayloadScheme{
		FieldConfigurationItems: []*models.FieldConfigurationItemScheme{
			{
				ID:          plan.Item.ID.Value,
				IsHidden:    plan.Item.IsHidden.Value,
				IsRequired:  plan.Item.IsRequired.Value,
				Description: plan.Item.Description.Value,
				Renderer:    plan.Item.Renderer.Value,
			},
		},
	}

	res, err := r.p.jira.Issue.Field.Configuration.Item.Update(ctx, issueFieldConfigurationId, &createRequestPayload)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue field configuration item, got error: %s\n%s", err, resBody))
		return
	}

	items, res, err := r.p.jira.Issue.Field.Configuration.Item.Gets(ctx, issueFieldConfigurationId, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration items, got error: %s\n%s", err, resBody))
		return
	}

	for _, i := range items.Values {
		if i.ID == plan.Item.ID.Value {
			plan.Item = &jiraIssueFieldConfigurationItem{
				ID:          types.String{Value: plan.Item.ID.Value},
				Description: types.String{Value: i.Description},
				IsHidden:    types.Bool{Value: i.IsHidden},
				IsRequired:  types.Bool{Value: i.IsRequired},
				Renderer:    types.String{Value: i.Renderer},
			}
		}
	}
	tflog.Debug(ctx, "Created issue field configuration item")

	plan.ID = types.String{Value: createIssueFieldConfigurationItemID(plan.IssueFieldConfiguration.Value, plan.Item.ID.Value)}

	tflog.Debug(ctx, "Storing issue field configuration item info into the state", map[string]interface{}{
		"issueFieldConfigurationItem": fmt.Sprintf("%+v, %+v", plan, *plan.Item),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationItemResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Reading issue field configuration item")

	var state jiraIssueFieldConfigurationItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration item from state", map[string]interface{}{
		"issueFieldConfigurationItem": fmt.Sprintf("%+v, %+v", state, *state.Item),
	})

	issueFieldConfigurationId, _ := strconv.Atoi(state.IssueFieldConfiguration.Value)
	issueFieldConfigurationItem, res, err := r.p.jira.Issue.Field.Configuration.Item.Gets(ctx, issueFieldConfigurationId, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration item, got error: %s\n%s", err, resBody))
		return
	}

	for _, i := range issueFieldConfigurationItem.Values {
		if i.ID == state.Item.ID.Value {
			state.Item = &jiraIssueFieldConfigurationItem{
				ID:          types.String{Value: state.Item.ID.Value},
				Description: types.String{Value: i.Description},
				IsHidden:    types.Bool{Value: i.IsHidden},
				IsRequired:  types.Bool{Value: i.IsRequired},
				Renderer:    types.String{Value: i.Renderer},
			}
		}
	}
	tflog.Debug(ctx, "Retrieved issue field configuration item from API state")

	state.ID = types.String{Value: createIssueFieldConfigurationItemID(state.IssueFieldConfiguration.Value, state.Item.ID.Value)}

	tflog.Debug(ctx, "Storing issue field configuration item into the state", map[string]interface{}{
		"issueFieldConfigurationItem": fmt.Sprintf("%+v, %+v", state, *state.Item),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueFieldConfigurationItemResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating issue field configuration item")

	var plan jiraIssueFieldConfigurationItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration item plan", map[string]interface{}{
		"issueFieldConfigurationItemPlan": fmt.Sprintf("%+v, %+v", plan, *plan.Item),
	})

	var state jiraIssueFieldConfigurationItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration item from state", map[string]interface{}{
		"issueFieldConfigurationItemState": fmt.Sprintf("%+v, %+v", state, *state.Item),
	})

	updateRequestPayload := models.UpdateFieldConfigurationItemPayloadScheme{
		FieldConfigurationItems: []*models.FieldConfigurationItemScheme{
			{
				ID:          plan.Item.ID.Value,
				IsHidden:    plan.Item.IsHidden.Value,
				IsRequired:  plan.Item.IsRequired.Value,
				Description: plan.Item.Description.Value,
				Renderer:    plan.Item.Renderer.Value,
			},
		},
	}

	issueFieldConfigurationId, _ := strconv.Atoi(plan.IssueFieldConfiguration.Value)
	res, err := r.p.jira.Issue.Field.Configuration.Item.Update(ctx, issueFieldConfigurationId, &updateRequestPayload)
	if err != nil {
		var resBody string
		if err != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue field configuration item, got error: %s\n%s", err, resBody))
		return
	}

	items, res, err := r.p.jira.Issue.Field.Configuration.Item.Gets(ctx, issueFieldConfigurationId, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration items, got error: %s\n%s", err, resBody))
		return
	}

	for _, i := range items.Values {
		if i.ID == plan.Item.ID.Value {
			plan.Item = &jiraIssueFieldConfigurationItem{
				ID:          types.String{Value: plan.Item.ID.Value},
				Description: types.String{Value: i.Description},
				IsHidden:    types.Bool{Value: i.IsHidden},
				IsRequired:  types.Bool{Value: i.IsRequired},
				Renderer:    types.String{Value: i.Renderer},
			}
		}
	}

	tflog.Debug(ctx, "Updated issue field configuration item in API state")

	plan.ID = types.String{Value: state.ID.Value}

	tflog.Debug(ctx, "Storing issue field configuration item plan into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationItemResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Warn(ctx, "Cannot destroy atlassian_jira_issue_field_configuration_item resource. Terraform will only remove this resource from the state file.")
	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}

func createIssueFieldConfigurationItemID(issueFieldConfiguration, item string) string {
	return strings.Join([]string{issueFieldConfiguration, item}, "-")
}

func (r *jiraIssueFieldConfigurationItemResource) checkIssueFieldConfigurationItemRenderable(ctx context.Context, p *jiraIssueFieldConfigurationItemResourceModel) error {
	var isRenderable bool
	searchPayload := models.FieldSearchOptionsScheme{
		IDs:    []string{p.Item.ID.Value},
		Expand: []string{"isLocked"},
	}

	itemDetails, res, err := r.p.jira.Issue.Field.Search(ctx, &searchPayload, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		return fmt.Errorf(" Unable to find issue field configuration item, got error: %s\n%s", err, resBody)
	}
	tflog.Debug(ctx, "Found issue field configuration item details", map[string]interface{}{
		"issueFieldConfigurationItem": fmt.Sprintf("%+v, %+v", itemDetails.Values[0], itemDetails.Values[0].Schema),
	})

	if itemDetails.Values[0].ID != p.Item.ID.Value {
		return fmt.Errorf(" Search result does not match issue field configuration item with ID: [%s]", p.Item.ID.Value)
	}

	if itemDetails.Values[0].IsLocked {
		return fmt.Errorf(" Tried to set a renderer for the locked item with ID: [%s]", p.Item.ID.Value)
	}

	isRenderable = strings.Contains(strings.Join(renderableItemTypes, ","), itemDetails.Values[0].Schema.Type)
	if !isRenderable {
		return fmt.Errorf(" Tried to set a renderer for the non-renderable item with ID: [%s]", p.Item.ID.Value)
	}

	return nil
}
