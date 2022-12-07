package atlassian

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueTypeDataSource struct {
		p atlassianProvider
	}

	jiraIssueTypeDataSourceModel struct {
		ID             types.String `tfsdk:"id"`
		Name           types.String `tfsdk:"name"`
		Description    types.String `tfsdk:"description"`
		HierarchyLevel types.Int64  `tfsdk:"hierarchy_level"`
		IconURL        types.String `tfsdk:"icon_url"`
		AvatarID       types.Int64  `tfsdk:"avatar_id"`
	}
)

var (
	_ datasource.DataSource = (*jiraIssueTypeDataSource)(nil)
)

func NewJiraIssueTypeDataSource() datasource.DataSource {
	return &jiraIssueTypeDataSource{}
}

func (*jiraIssueTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_type"
}

func (*jiraIssueTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Jira Issue Type Data Source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the issue type.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the issue type.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the issue type.",
				Computed:            true,
			},
			"hierarchy_level": schema.Int64Attribute{
				MarkdownDescription: "The hierarchy level of the issue type.",
				Computed:            true,
			},
			"icon_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the issue type's avatar.",
				Computed:            true,
			},
			"avatar_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the issue type's avatar.",
				Computed:            true,
			},
		},
	}
}

func (d *jiraIssueTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraIssueTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type data source")

	var newstate jiraIssueTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newstate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	issueType, res, err := d.p.jira.Issue.Type.Get(ctx, newstate.ID.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue type from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v", issueType),
	})

	newstate.Name = types.StringValue(issueType.Name)
	newstate.Description = types.StringValue(issueType.Description)
	newstate.HierarchyLevel = types.Int64Value(int64(issueType.HierarchyLevel))
	newstate.IconURL = types.StringValue(issueType.IconURL)
	newstate.AvatarID = types.Int64Value(int64(issueType.AvatarID))

	tflog.Debug(ctx, "Storing issue type into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newstate)...)
}
