package atlassian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueFieldConfigurationDataSource struct {
		p atlassianProvider
	}

	jiraIssueFieldConfigurationDataSourceType struct{}

	jiraIssueFieldConfigurationDataSourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ datasource.DataSource   = (*jiraIssueFieldConfigurationDataSource)(nil)
	_ provider.DataSourceType = (*jiraIssueFieldConfigurationDataSourceType)(nil)
)

func (d *jiraIssueFieldConfigurationDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Field Configuration Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue field configuration.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue field configuration.",
				Computed:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the issue field configuration.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *jiraIssueFieldConfigurationDataSourceType) NewDataSource(_ context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraIssueFieldConfigurationDataSource{
		p: provider,
	}, diags

}

func (d *jiraIssueFieldConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue field configuration")

	var newState jiraIssueFieldConfigurationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration config", map[string]interface{}{
		"issueFieldConfiguration": fmt.Sprintf("%+v", newState),
	})

	issueFieldConfigurationId, err := strconv.Atoi(newState.ID.Value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Unable to parse value of \"id\" attribute.", "Value of \"id\" attribute can only be a numeric string.")
		return
	}

	issueFieldConfiguration, res, err := d.p.jira.Issue.Field.Configuration.Gets(ctx, []int{issueFieldConfigurationId}, false, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue field configuration, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue field configuration from API state", map[string]interface{}{
		"issueFieldConfiguration": fmt.Sprintf("%+v", issueFieldConfiguration.Values[0]),
	})

	newState.Name = types.String{Value: issueFieldConfiguration.Values[0].Name}
	newState.Description = types.String{Value: issueFieldConfiguration.Values[0].Description}

	tflog.Debug(ctx, "Storing issue field configuration into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
