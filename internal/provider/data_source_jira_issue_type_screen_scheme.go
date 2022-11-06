package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraIssueTypeScreenSchemeDataSource struct {
		p atlassianProvider
	}

	jiraIssueTypeScreenSchemeDataSourceModel struct {
		ID                types.String                       `tfsdk:"id"`
		Name              types.String                       `tfsdk:"name"`
		Description       types.String                       `tfsdk:"description"`
		IssueTypeMappings []jiraIssueTypeScreenSchemeMapping `tfsdk:"issue_type_mappings"`
	}
)

var (
	_ datasource.DataSource = (*jiraIssueTypeScreenSchemeDataSource)(nil)
)

func NewJiraIssueTypeScreenSchemeDataSource() datasource.DataSource {
	return &jiraIssueTypeScreenSchemeDataSource{}
}

func (*jiraIssueTypeScreenSchemeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_issue_type_screen_scheme"
}

func (*jiraIssueTypeScreenSchemeDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Issue Type Screen Scheme Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the issue type screen scheme.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "The name of the issue type screen scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Computed: true,
				Type:     types.StringType,
			},
			"description": {
				MarkdownDescription: "The description of the issue type screen scheme. " +
					"The maximum length is 255 characters.",
				Computed: true,
				Type:     types.StringType,
			},
			"issue_type_mappings": {
				MarkdownDescription: "The IDs of the screen schemes for the issue type IDs and default. " +
					"A default entry is required to create an issue type screen scheme, it defines the mapping for all issue types without a screen scheme.",
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"issue_type_id": {
							MarkdownDescription: "The ID of the issue type or default. " +
								"Only issue types used in classic projects are accepted. " +
								"An entry for default must be provided and defines the mapping for all issue types without a screen scheme.",
							Computed: true,
							Type:     types.StringType,
						},
						"screen_scheme_id": {
							MarkdownDescription: "The ID of the screen scheme. " +
								"Only screen schemes used in classic projects are accepted.",
							Computed: true,
							Type:     types.StringType,
						},
					},
				),
			},
		},
	}, nil
}

func (d *jiraIssueTypeScreenSchemeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraIssueTypeScreenSchemeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading issue type screen scheme data source")

	var newState jiraIssueTypeScreenSchemeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded issue type screen scheme config", map[string]interface{}{
		"readConfig": fmt.Sprintf("%+v", newState),
	})

	issueTypeScreenSchemeId, err := strconv.Atoi(newState.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Unable to parse value of \"id\" attribute.", "Value of \"id\" attribute can only be a numeric string.")
		return
	}
	options := &models.ScreenSchemeParamsScheme{
		IDs: []int{issueTypeScreenSchemeId},
	}

	issueTypeScreenScheme, res, err := d.p.jira.Issue.Type.ScreenScheme.Gets(ctx, options, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type screen scheme, got error: %s\n%s", err, resBody))
		return
	}

	issueTypeMappings, res, err := d.p.jira.Issue.Type.ScreenScheme.Mapping(ctx, []int{issueTypeScreenSchemeId}, 0, 50)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get issue type screen scheme mappings, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved issue type screen scheme from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v, Mappings:%+v", issueTypeScreenScheme.Values[0], issueTypeMappings.Values[0]),
	})

	newState.Name = types.String{Value: issueTypeScreenScheme.Values[0].Name}
	newState.Description = types.String{Value: issueTypeScreenScheme.Values[0].Description}
	var mappings []jiraIssueTypeScreenSchemeMapping
	for _, m := range issueTypeMappings.Values {
		mappings = append(mappings, jiraIssueTypeScreenSchemeMapping{
			IssueTypeId:    types.String{Value: m.IssueTypeID},
			ScreenSchemeId: types.String{Value: m.ScreenSchemeID},
		})
	}
	newState.IssueTypeMappings = mappings

	tflog.Debug(ctx, "Storing issue type screen scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
