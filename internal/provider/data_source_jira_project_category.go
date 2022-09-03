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
	jiraProjectCategoryDataSource struct {
		p atlassianProvider
	}

	jiraProjectCategoryDataSourceType struct{}

	jiraProjectCategoryDataSourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
		Self        types.String `tfsdk:"self"`
	}
)

var (
	_ datasource.DataSource   = (*jiraProjectCategoryDataSource)(nil)
	_ provider.DataSourceType = (*jiraProjectCategoryDataSourceType)(nil)
)

func (d *jiraProjectCategoryDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (d *jiraProjectCategoryDataSourceType) NewDataSource(_ context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraProjectCategoryDataSource{
		p: provider,
	}, diags
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

	projectCategoryId, err := strconv.Atoi(newState.ID.Value)
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

	newState.Name = types.String{Value: projectCategory.Name}
	newState.Description = types.String{Value: projectCategory.Description}
	newState.Self = types.String{Value: projectCategory.Self}

	tflog.Debug(ctx, "Storing project category into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
