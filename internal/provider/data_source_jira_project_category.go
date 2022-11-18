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
	jiraProjectCategoryDataSource struct {
		p atlassianProvider
	}

	jiraProjectCategoryDataSourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
		Self        types.String `tfsdk:"self"`
	}
)

var (
	_ datasource.DataSource = (*jiraProjectCategoryDataSource)(nil)
)

func NewJiraProjectCategoryDataSource() datasource.DataSource {
	return &jiraProjectCategoryDataSource{}
}

func (*jiraProjectCategoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_project_category"
}

func (*jiraProjectCategoryDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Project Category Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the project category.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the project category.",
				Computed:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the project category.",
				Computed:            true,
				Type:                types.StringType,
			},
			"self": {
				MarkdownDescription: "The URL of the project category.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *jiraProjectCategoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraProjectCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading project category data source")

	var newState jiraProjectCategoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded project category config", map[string]interface{}{
		"readConfig": fmt.Sprintf("%+v", newState),
	})

	projectCategoryId, err := strconv.Atoi(newState.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Unable to parse value of \"id\" attribute.", "Value of \"id\" attribute can only be a numeric string.")
		return
	}

	projectCategory, res, err := d.p.jira.Project.Category.Get(ctx, projectCategoryId)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get project category, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved project category from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v", projectCategory),
	})

	newState.Name = types.StringValue(projectCategory.Name)
	newState.Description = types.StringValue(projectCategory.Description)
	newState.Self = types.StringValue(projectCategory.Self)

	tflog.Debug(ctx, "Storing project category into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
