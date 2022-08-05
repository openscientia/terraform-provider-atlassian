package atlassian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueTypeScreenSchemeDataSource struct {
		p provider
	}

	jiraIssueTypeScreenSchemeDataSourceType struct{}

	jiraIssueTypeScreenSchemeDataSourceModel struct {
		ID                types.String                       `tfsdk:"id"`
		Name              types.String                       `tfsdk:"name"`
		Description       types.String                       `tfsdk:"description"`
		IssueTypeMappings []jiraIssueTypeScreenSchemeMapping `tfsdk:"issue_type_mappings"`
	}
)

var (
	_ tfsdk.DataSource     = (*jiraIssueTypeScreenSchemeDataSource)(nil)
	_ tfsdk.DataSourceType = (*jiraIssueTypeScreenSchemeDataSourceType)(nil)
)

func (d *jiraIssueTypeScreenSchemeDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Type Screen Scheme Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type screen scheme.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue type screen scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Computed: true,
				Type:     types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the issue type screen scheme. " +
					"The maximum length is 255 characters.",
				Computed: true,
				Type:     types.StringType,
			},
			"issue_type_mappings": {
				MarkdownDescription: "The IDs of the screen schemes for the issue type IDs and default. " +
					"A default entry is required to create an issue type screen scheme, it defines the mapping for all issue types without a screen scheme.",
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"issue_type_id": {
							MarkdownDescription: "The ID of the issue type or default. " +
								"Only issue types used in classic projects are accepted. " +
								"An entry for default must be provided and defines the mapping for all issue types without a screen scheme.",
							Computed: true,
							Type:     types.StringType,
						},
						"screen_scheme_id": {
							MarkdownDescription: "The ID of the screen scheme. " +
								"Only screen schemes used in classic projects are accepted.",
							Computed: true,
							Type:     types.StringType,
						},
					},
				),
			},
		},
	}, nil
}

func (d *jiraIssueTypeScreenSchemeDataSourceType) NewDataSource(_ context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraIssueTypeScreenSchemeDataSource{
		p: provider,
	}, diags
}

func (d *jiraIssueTypeScreenSchemeDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Reading issue type screen scheme")

	var newState jiraIssueTypeScreenSchemeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme config", map[string]interface{}{
		"issueTypeScreenScheme": fmt.Sprintf("%+v", newState),
	})

	issueTypeScreenSchemeId, err := strconv.Atoi(newState.ID.Value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Unable to parse value of \"id\" attribute.", "Value of \"id\" attribute can only be a numeric string.")
		return
	}

	issueTypeScreenScheme, res, err := d.p.jira.Issue.Type.ScreenScheme.Gets(ctx, []int{issueTypeScreenSchemeId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type screen scheme, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue type screen scheme from API state", map[string]interface{}{
		"issueTypeScreenScheme": fmt.Sprintf("%+v", issueTypeScreenScheme.Values[0]),
	})

	newState.Name = types.String{Value: issueTypeScreenScheme.Values[0].Name}
	newState.Description = types.String{Value: issueTypeScreenScheme.Values[0].Description}

	issueTypeMappings, res, err := d.p.jira.Issue.Type.ScreenScheme.Mapping(ctx, []int{issueTypeScreenSchemeId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type screen scheme mappings, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue type screen scheme mappings from API state", map[string]interface{}{
		"issueTypeMappings": fmt.Sprintf("%+v", issueTypeMappings.Values[0]),
	})
	var mappings []jiraIssueTypeScreenSchemeMapping
	for _, m := range issueTypeMappings.Values {
		mappings = append(mappings, jiraIssueTypeScreenSchemeMapping{
			IssueTypeId:    types.String{Value: m.IssueTypeID},
			ScreenSchemeId: types.String{Value: m.ScreenSchemeID},
		})
	}
	newState.IssueTypeMappings = mappings

	tflog.Debug(ctx, "Storing issue type screen scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
