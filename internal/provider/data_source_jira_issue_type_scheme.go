package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueTypeSchemeDataSource struct {
		p atlassianProvider
	}

	jiraIssueTypeSchemeDataSourceModel struct {
		ID                 types.String `tfsdk:"id"`
		Name               types.String `tfsdk:"name"`
		Description        types.String `tfsdk:"description"`
		DefaultIssueTypeId types.String `tfsdk:"default_issue_type_id"`
		IssueTypeIds       types.List   `tfsdk:"issue_type_ids"`
	}
)

var (
	_ datasource.DataSource = (*jiraIssueTypeSchemeDataSource)(nil)
)

func NewJiraIssueTypeSchemeDataSource() datasource.DataSource {
	return &jiraIssueTypeSchemeDataSource{}
}

func (*jiraIssueTypeSchemeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_type_scheme"
}

func (*jiraIssueTypeSchemeDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Type Scheme Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type scheme.",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "The name of the issue type scheme.",
				Type:                types.StringType,
				Computed:            true,
			},
			"description": {
				MarkdownDescription: "The description of the issue type scheme.",
				Type:                types.StringType,
				Computed:            true,
			},
			"default_issue_type_id": {
				MarkdownDescription: "The ID of the default issue type of the issue type scheme.",
				Type:                types.StringType,
				Computed:            true,
			},
			"issue_type_ids": {
				MarkdownDescription: "The list of issue types IDs of the issue type scheme.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed: true,
			},
		},
	}, nil
}

func (d *jiraIssueTypeSchemeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraIssueTypeSchemeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type scheme data source")

	var newState jiraIssueTypeSchemeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	issueTypeSchemeID, err := strconv.Atoi(newState.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Conversion failed: %s", err))
		return
	}

	// Get issue type scheme details
	issueTypeScheme, res, err := d.p.jira.Issue.Type.Scheme.Gets(ctx, []int{issueTypeSchemeID}, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type, got error: %s\n%s", err, resBody))
		return
	}

	// Get issue type scheme items
	issueTypeSchemeItems, res, err := d.p.jira.Issue.Type.Scheme.Items(ctx, []int{issueTypeSchemeID}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type scheme items, got error: %s\n%s", err, resBody))
		return
	}

	tflog.Debug(ctx, "Retrieved issue type scheme from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v, items:%+v", issueTypeScheme.Values[0], issueTypeSchemeItems.Values),
	})

	ids := types.ListNull(types.StringType)
	for _, item := range issueTypeSchemeItems.Values {
		id := types.StringValue(item.IssueTypeID)
		ids, _ = types.ListValue(types.StringType, append(ids.Elements(), id))
	}

	newState.ID = types.String{Value: issueTypeScheme.Values[0].ID}
	newState.Name = types.String{Value: issueTypeScheme.Values[0].Name}
	newState.Description = types.String{Value: issueTypeScheme.Values[0].Description}
	newState.DefaultIssueTypeId = types.String{Value: issueTypeScheme.Values[0].DefaultIssueTypeID}
	newState.IssueTypeIds = ids

	tflog.Debug(ctx, "Storing issue type scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
