package atlassian

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.DataSourceType = jiraIssueTypeDataSourceType{}
var _ tfsdk.DataSource = jiraIssueTypeDataSource{}

type jiraIssueTypeDataSourceType struct{}

type jiraIssueTypeDataSourceData struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	HierarchyLevel types.Int64  `tfsdk:"hierarchy_level"`
	IconURL        types.String `tfsdk:"icon_url"`
	AvatarID       types.Int64  `tfsdk:"avatar_id"`
}

type jiraIssueTypeDataSource struct {
	p provider
}

func (t jiraIssueTypeDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
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
		Version: 1,
	}, nil
}

func (t jiraIssueTypeDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jiraIssueTypeDataSource{
		p: provider,
	}, diags
}

func (d jiraIssueTypeDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data jiraIssueTypeDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	returnedIssueType, res, err := d.p.jira.Issue.Type.Get(ctx, data.ID.Value)
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
