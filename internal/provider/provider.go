package atlassian

import (
	"context"
	"os"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/validators"
)

type (
	atlassianProvider struct {
		jira *jira.Client

		version string
	}

	atlassianProviderModel struct {
		Url      types.String `tfsdk:"url"`
		Username types.String `tfsdk:"username"`
		ApiToken types.String `tfsdk:"apitoken"`
	}
)

var (
	_ provider.Provider             = (*atlassianProvider)(nil)
	_ provider.ProviderWithMetadata = (*atlassianProvider)(nil)
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &atlassianProvider{
			version: version,
		}
	}
}

func (p *atlassianProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "atlassian"
}

func (*atlassianProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Atlassian Provider",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Atlassian Host URL. Can also be set with the `ATLASSIAN_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					validators.UrlWithScheme("https"),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Atlassian Username. Can also be set with the `ATLASSIAN_USERNAME` environment variable.",
				Optional:            true,
			},
			"apitoken": schema.StringAttribute{
				MarkdownDescription: "Atlassian API Token. Can also be set with the `ATLASSIAN_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *atlassianProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data atlassianProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a user to the provider
	var username string
	if data.Username.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client.",
			"Cannot use unknown value as Username",
		)
		return
	}
	if data.Username.IsNull() {
		username = os.Getenv("ATLASSIAN_USERNAME")
	} else {
		username = data.Username.ValueString()
	}
	if username == "" {
		resp.Diagnostics.AddError(
			"Unable to find Username value.",
			"Username cannot be an empty string.",
		)
		return
	}

	// User must provide a password to the provider
	var apitoken string
	if data.ApiToken.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client.",
			"Cannot use unknown value as ApiToken.",
		)
		return
	}

	if data.ApiToken.IsNull() {
		apitoken = os.Getenv("ATLASSIAN_TOKEN")
	} else {
		apitoken = data.ApiToken.ValueString()
	}

	if apitoken == "" {
		resp.Diagnostics.AddError(
			"Unable to find ApiToken.",
			"ApiToken cannot be an empty string.",
		)
		return
	}

	// User must specify a host
	var url string
	if data.Url.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client.",
			"Cannot use unknown value as Url.",
		)
		return
	}

	if data.Url.IsNull() {
		url = os.Getenv("ATLASSIAN_URL")
	} else {
		url = data.Url.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError(
			"Unable to find Url.",
			"Url cannot be an empty string.",
		)
		return
	}

	c, err := jira.New(nil, url)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create Atlassian client:\n\n"+err.Error(),
		)
		return
	}
	c.Auth.SetBasicAuth(username, apitoken)

	p.jira = c

	resp.DataSourceData = p.jira
	resp.ResourceData = p.jira
}

func (*atlassianProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJiraGroupResource,
		NewJiraGroupUserResource,
		NewJiraIssueFieldConfigurationItemResource,
		NewJiraIssueFieldConfigurationResource,
		NewJiraIssueFieldConfigurationSchemeMappingResource,
		NewJiraIssueFieldConfigurationSchemeResource,
		NewJiraIssueScreenResource,
		NewJiraIssueTypeResource,
		NewJiraIssueTypeSchemeResource,
		NewJiraIssueTypeScreenSchemeResource,
		NewJiraPermissionGrantResource,
		NewJiraPermissionSchemeResource,
		NewJiraProjectCategoryResource,
		NewJiraScreenSchemeResource,
	}
}

func (*atlassianProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewJiraGroupDataSource,
		NewJiraIssueFieldConfigurationDataSource,
		NewJiraIssueFieldConfigurationSchemeDataSource,
		NewJiraIssueScreenDataSource,
		NewJiraIssueTypeDataSource,
		NewJiraIssueTypeSchemeDataSource,
		NewJiraIssueTypeScreenSchemeDataSource,
		NewJiraMyselfDataSource,
		NewJiraPermissionGrantDataSource,
		NewJiraPermissionSchemeDataSource,
		NewJiraProjectCategoryDataSource,
		NewJiraScreenSchemeDataSource,
		NewJiraServerInfoDataSource,
	}
}
