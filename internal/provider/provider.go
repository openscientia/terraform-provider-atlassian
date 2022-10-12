package atlassian

import (
	"context"
	"os"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_validation"
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

func (*atlassianProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Atlassian Provider",
		Attributes: map[string]tfsdk.Attribute{
			"url": {
				MarkdownDescription: "Atlassian Host URL. Can also be set with the `ATLASSIAN_URL` environment variable.",
				Computed:            true,
				Optional:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					attribute_validation.UrlWithScheme("https"),
				},
			},
			"username": {
				MarkdownDescription: "Atlassian Username. Can also be set with the `ATLASSIAN_USERNAME` environment variable.",
				Computed:            true,
				Optional:            true,
				Type:                types.StringType,
			},
			"apitoken": {
				MarkdownDescription: "Atlassian API Token. Can also be set with the `ATLASSIAN_TOKEN` environment variable.",
				Computed:            true,
				Optional:            true,
				Sensitive:           true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (p *atlassianProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data atlassianProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a user to the provider
	var username string
	if data.Username.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client.",
			"Cannot use unknown value as Username",
		)
		return
	}
	if data.Username.Null {
		username = os.Getenv("ATLASSIAN_USERNAME")
	} else {
		username = data.Username.Value
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
	if data.ApiToken.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client.",
			"Cannot use unknown value as ApiToken.",
		)
		return
	}

	if data.ApiToken.Null {
		apitoken = os.Getenv("ATLASSIAN_TOKEN")
	} else {
		apitoken = data.ApiToken.Value
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
	if data.Url.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client.",
			"Cannot use unknown value as Url.",
		)
		return
	}

	if data.Url.Null {
		url = os.Getenv("ATLASSIAN_URL")
	} else {
		url = data.Url.Value
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
