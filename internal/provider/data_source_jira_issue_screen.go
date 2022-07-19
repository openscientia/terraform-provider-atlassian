package atlassian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueScreenDataSourceType struct{}
	jiraIssueScreenDataSource     struct {
		p provider
	}
	jiraIssueScreenDataSourceModel struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
	}
)

var (
	_ tfsdk.DataSourceType = jiraIssueScreenDataSourceType{}
	_ tfsdk.DataSource     = jiraIssueScreenDataSource{}
)

func (jiraIssueScreenDataSourceType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Screen Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue screen.",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "The name of the screen." +
					"The name must be unique." +
					"The maximum length is 255 characters.",
				Type:     types.StringType,
				Computed: true,
			},
			"description": {
				MarkdownDescription: "The description of the screen." +
					"The maximum length is 255 characters.",
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

func (d jiraIssueScreenDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jiraIssueScreenDataSource{
		p: provider,
	}, diags
}

func (d jiraIssueScreenDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Reading issue screen")
	var newState jiraIssueScreenDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue screen config", map[string]interface{}{
		"issueScreen": fmt.Sprintf("%+v", newState),
	})

	issueScreenId, _ := strconv.Atoi(newState.ID.Value)

	issueScreen, res, err := d.p.jira.Screen.Gets(ctx, []int{issueScreenId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue screen, got error: %s\n%s", err.Error(), resBody))
		return
	}
	tflog.Debug(ctx, "Retrieve issue screen from API state")

	newState.Name = types.String{Value: issueScreen.Values[0].Name}
	newState.Description = types.String{Value: issueScreen.Values[0].Description}

	tflog.Debug(ctx, "Storing issue screen info into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
