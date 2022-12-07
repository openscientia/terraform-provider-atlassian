package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraScreenSchemeDataSource struct {
		p atlassianProvider
	}

	jiraScreenSchemeDataSourceModel struct {
		ID          types.String                `tfsdk:"id"`
		Name        types.String                `tfsdk:"name"`
		Description types.String                `tfsdk:"description"`
		Screens     *jiraScreenSchemeTypesModel `tfsdk:"screens"`
	}
)

var (
	_ datasource.DataSource = (*jiraScreenSchemeDataSource)(nil)
)

func NewJiraScreenSchemeDataSource() datasource.DataSource {
	return &jiraScreenSchemeDataSource{}
}

func (*jiraScreenSchemeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_screen_scheme"
}

func (d *jiraScreenSchemeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Jira Screen Scheme Data Source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the screen scheme.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the screen scheme. " +
					"The name must be unique. " +
					"The maximum length is 255 characters.",
				Computed: true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the screen scheme. " +
					"The maximum length is 255 characters.",
				Computed: true,
			},
			"screens": schema.SingleNestedAttribute{
				MarkdownDescription: "The IDs of the screens for the screen types of the screen scheme. " +
					"Only screens used in classic projects are accepted.",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"create": schema.Int64Attribute{
						MarkdownDescription: "The ID of the create screen.",
						Computed:            true,
					},
					"default": schema.Int64Attribute{
						MarkdownDescription: "The ID of the default screen. Required when creating a screen scheme.",
						Computed:            true,
					},
					"view": schema.Int64Attribute{
						MarkdownDescription: "The ID of the view screen.",
						Computed:            true,
					},
					"edit": schema.Int64Attribute{
						MarkdownDescription: "The ID of the edit screen.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *jiraScreenSchemeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraScreenSchemeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading screen scheme data source")

	var newState jiraScreenSchemeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded screen scheme config", map[string]interface{}{
		"readConfig": fmt.Sprintf("%+v", newState),
	})

	screenSchemeId, err := strconv.Atoi(newState.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Unable to parse value of \"id\" attribute.", "Value of \"id\" attribute can only be a numeric string.")
		return
	}

	options := &models.ScreenSchemeParamsScheme{
		IDs: []int{screenSchemeId},
	}
	screenScheme, res, err := d.p.jira.Screen.Scheme.Gets(ctx, options, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get screen scheme, got error: %s\n%s", err, resBody))
	}
	tflog.Debug(ctx, "Retrieved screen scheme from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v", screenScheme.Values[0]),
	})

	newState.Name = types.StringValue(screenScheme.Values[0].Name)
	newState.Description = types.StringValue(screenScheme.Values[0].Description)
	newState.Screens = &jiraScreenSchemeTypesModel{
		Create:  types.Int64Value(int64(screenScheme.Values[0].Screens.Create)),
		Default: types.Int64Value(int64(screenScheme.Values[0].Screens.Default)),
		View:    types.Int64Value(int64(screenScheme.Values[0].Screens.View)),
		Edit:    types.Int64Value(int64(screenScheme.Values[0].Screens.Edit)),
	}

	tflog.Debug(ctx, "Storing screen scheme into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
