package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueFieldConfigurationSchemeDataSource struct {
		p atlassianProvider
	}

	jiraIssueFieldConfigurationSchemeDataSourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ datasource.DataSource = (*jiraIssueFieldConfigurationSchemeDataSource)(nil)
)

func NewJiraIssueFieldConfigurationSchemeDataSource() datasource.DataSource {
	return &jiraIssueFieldConfigurationSchemeDataSource{}
}

func (*jiraIssueFieldConfigurationSchemeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_field_configuration_scheme"
}

func (*jiraIssueFieldConfigurationSchemeDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Field Configuration Scheme Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue field configuration scheme.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue field configuration scheme.",
				Computed:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the issue field configuration scheme.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *jiraIssueFieldConfigurationSchemeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jira.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jira.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.p.jira = client
}

func (d *jiraIssueFieldConfigurationSchemeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue field configuration scheme")

	var newState jiraIssueFieldConfigurationSchemeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue field configuration scheme config", map[string]interface{}{
		"issueFieldConfiguration": fmt.Sprintf("%+v", newState),
	})

	issueFieldConfigurationSchemeId, err := strconv.Atoi(newState.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Unable to parse value of \"id\" attribute.", "Value of \"id\" attribute can only be a numeric string.")
		return
	}

	issueFieldConfigurationScheme, res, err := d.p.jira.Issue.Field.Configuration.Scheme.Gets(ctx, []int{issueFieldConfigurationSchemeId}, 0, 1)
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

	newState.Name = types.String{Value: issueFieldConfigurationScheme.Values[0].Name}
	newState.Description = types.String{Value: issueFieldConfigurationScheme.Values[0].Description}

	tflog.Debug(ctx, "Storing issue field configuration scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
