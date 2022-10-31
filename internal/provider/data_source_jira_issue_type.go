package atlassian

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (*jiraIssueTypeDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Type Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type.",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "The name of the issue type.",
				Type:                types.StringType,
				Computed:            true,
			},
			"description": {
				MarkdownDescription: "The description of the issue type.",
				Type:                types.StringType,
				Computed:            true,
			},
			"hierarchy_level": {
				MarkdownDescription: "The hierarchy level of the issue type.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"icon_url": {
				MarkdownDescription: "The URL of the issue type's avatar.",
				Type:                types.StringType,
				Computed:            true,
			},
			"avatar_id": {
				MarkdownDescription: "The ID of the issue type's avatar.",
				Type:                types.Int64Type,
				Computed:            true,
			},
		},
	}, nil
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
	var data jiraIssueTypeDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	returnedIssueType, res, err := d.p.jira.Issue.Type.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type, got error: %s\n%s", err.Error(), res.Bytes.String()))
		return
	}

	data.ID = types.String{Value: returnedIssueType.ID}
	data.Name = types.String{Value: returnedIssueType.Name}
	data.Description = types.String{Value: returnedIssueType.Description}
	data.HierarchyLevel = types.Int64{Value: int64(returnedIssueType.HierarchyLevel)}
	data.IconURL = types.String{Value: returnedIssueType.IconURL}
	data.AvatarID = types.Int64{Value: int64(returnedIssueType.AvatarID)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
