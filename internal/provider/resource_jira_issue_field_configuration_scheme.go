package atlassian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_plan_modification"
)

type (
	jiraIssueFieldConfigurationSchemeResource struct {
		p provider
	}

	jiraIssueFieldConfigurationSchemeResourceType struct{}

	jiraIssueFieldConfigurationSchemeResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ tfsdk.Resource                = (*jiraIssueFieldConfigurationSchemeResource)(nil)
	_ tfsdk.ResourceType            = (*jiraIssueFieldConfigurationSchemeResourceType)(nil)
	_ tfsdk.ResourceWithImportState = (*jiraIssueFieldConfigurationSchemeResource)(nil)
)

func (*jiraIssueFieldConfigurationSchemeResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Field Configuration Scheme Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue field configuration scheme.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue field configuration scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": {
				MarkdownDescription: "The description of the issue field configuration scheme. " +
					"The maximum length is 1024 characters.",
				Optional: true,
				Computed: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtMost(1024),
				},
				PlanModifiers: tfsdk.AttributePlanModifiers{
					attribute_plan_modification.DefaultValue(types.String{Value: ""}),
				},
			},
		},
	}, nil
}

func (r *jiraIssueFieldConfigurationSchemeResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraIssueFieldConfigurationSchemeResource{
		p: provider,
	}, diags
}

func (r *jiraIssueFieldConfigurationSchemeResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating issue field configuration scheme")

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
	}

	var plan jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	issueFieldConfigurationScheme, res, err := r.p.jira.Issue.Field.Configuration.Scheme.Create(ctx, plan.Name.Value, plan.Description.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create issue field configuration scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created issue field configuration scheme", map[string]interface{}{
		"issueFieldConfigurationScheme": fmt.Sprintf("%+v", issueFieldConfigurationScheme),
	})

	plan.ID = types.String{Value: issueFieldConfigurationScheme.ID}

	tflog.Debug(ctx, "Storing issue field configuration scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Reading issue field configuration scheme")

	var state jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	id, _ := strconv.Atoi(state.ID.Value)
	issueFieldConfigurationScheme, res, err := r.p.jira.Issue.Field.Configuration.Scheme.Gets(ctx, []int{id}, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue field configuration scheme from API state", map[string]interface{}{
		"issueFieldConfigurationScheme": fmt.Sprintf("%+v", issueFieldConfigurationScheme),
	})

	state.Name = types.String{Value: issueFieldConfigurationScheme.Values[0].Name}
	state.Description = types.String{Value: issueFieldConfigurationScheme.Values[0].Description}

	tflog.Debug(ctx, "Storing issue field configuration scheme info into the state", map[string]interface{}{
		"newState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating issue field configuration scheme")

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

	id, _ := strconv.Atoi(state.ID.Value)
	res, err := r.p.jira.Issue.Field.Configuration.Scheme.Update(ctx, id, plan.Name.Value, plan.Description.Value)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update issue field configuration scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Updated issue field configuration scheme")

	plan.ID = types.String{Value: state.ID.Value}

	tflog.Debug(ctx, "Storing issue field configuration scheme info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraIssueFieldConfigurationSchemeResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Deleting issue field configuration scheme")

	var state jiraIssueFieldConfigurationSchemeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme from state")

	id, _ := strconv.Atoi(state.ID.Value)
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
