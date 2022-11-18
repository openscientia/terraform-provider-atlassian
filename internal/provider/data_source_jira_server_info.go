package atlassian

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraServerInfoDataSource struct {
		p atlassianProvider
	}

	jiraServerInfoDataSourceModel struct {
		ID             types.String `tfsdk:"id"`
		BaseURL        types.String `tfsdk:"base_url"`
		Version        types.String `tfsdk:"version"`
		VersionNumbers types.List   `tfsdk:"version_numbers"`
		DeploymentType types.String `tfsdk:"deployment_type"`
		BuildNumber    types.Int64  `tfsdk:"build_number"`
		BuildDate      types.String `tfsdk:"build_date"`
		ServerTime     types.String `tfsdk:"server_time"`
		ScmInfo        types.String `tfsdk:"scm_info"`
		ServerTitle    types.String `tfsdk:"server_title"`
	}
)

var (
	_ datasource.DataSource = (*jiraServerInfoDataSource)(nil)
)

func NewJiraServerInfoDataSource() datasource.DataSource {
	return &jiraServerInfoDataSource{}
}

func (*jiraServerInfoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_server_info"
}

func (*jiraServerInfoDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Server Info Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of server info. Defaults to `base_url`.",
				Computed:            true,
				Type:                types.StringType,
			},
			"base_url": {
				MarkdownDescription: "The base URL of the Jira instance.",
				Computed:            true,
				Type:                types.StringType,
			},
			"version": {
				MarkdownDescription: "The version of Jira.",
				Computed:            true,
				Type:                types.StringType,
			},
			"version_numbers": {
				MarkdownDescription: "The major, minor, and revision version numbers of the Jira version.",
				Computed:            true,
				Type: types.ListType{
					ElemType: types.Int64Type,
				},
			},
			"deployment_type": {
				MarkdownDescription: "The type of server deployment. This is always returned as Cloud.",
				Computed:            true,
				Type:                types.StringType,
			},
			"build_number": {
				MarkdownDescription: "The build number of the Jira version.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"build_date": {
				MarkdownDescription: "The timestamp when the Jira version was built.",
				Computed:            true,
				Type:                types.StringType,
			},
			"server_time": {
				MarkdownDescription: "The time in Jira when this request was responded to.",
				Computed:            true,
				Type:                types.StringType,
			},
			"scm_info": {
				MarkdownDescription: "The unique identifier of the Jira version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"server_title": {
				MarkdownDescription: "The name of the Jira instance.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *jiraServerInfoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraServerInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading server info data source")

	serverInfo, res, err := d.p.jira.Server.Info(ctx)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get server info, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved server info from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v", serverInfo),
	})

	newState := &jiraServerInfoDataSourceModel{
		ID:             types.StringValue(serverInfo.BaseURL),
		BaseURL:        types.StringValue(serverInfo.BaseURL),
		Version:        types.StringValue(serverInfo.Version),
		VersionNumbers: types.ListNull(types.Int64Type),
		DeploymentType: types.StringValue(serverInfo.DeploymentType),
		BuildNumber:    types.Int64Value(int64(serverInfo.BuildNumber)),
		BuildDate:      types.StringValue(serverInfo.BuildDate),
		ServerTime:     types.StringValue(serverInfo.ServerTime),
		ScmInfo:        types.StringValue(serverInfo.ScmInfo),
		ServerTitle:    types.StringValue(serverInfo.ServerTitle),
	}
	newState.VersionNumbers, _ = types.ListValueFrom(ctx, types.Int64Type, serverInfo.VersionNumbers)

	tflog.Debug(ctx, "Storing server info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
