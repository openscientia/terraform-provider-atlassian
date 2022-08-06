package atlassian

import (
	"context"
	"fmt"
	"os"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/openscientia/terraform-provider-atlassian/internal/provider/attribute_validation"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.Provider = &provider{}

type provider struct {
	jira *jira.Client

	configured bool

	version string
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Atlassian",

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

		Version: 1,
	}, nil
}

type providerData struct {
	Url      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	ApiToken types.String `tfsdk:"apitoken"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

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
	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"atlassian_jira_issue_field_configuration": &jiraIssueFieldConfigurationResourceType{},
		"atlassian_jira_issue_screen":              jiraIssueScreenResourceType{},
		"atlassian_jira_issue_type":                jiraIssueTypeResourceType{},
		"atlassian_jira_issue_type_scheme":         jiraIssueTypeSchemeResourceType{},
		"atlassian_jira_issue_type_screen_scheme":  &jiraIssueTypeScreenSchemeResourceType{},
		"atlassian_jira_screen_scheme":             &jiraScreenSchemeResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"atlassian_jira_issue_screen":             jiraIssueScreenDataSourceType{},
		"atlassian_jira_issue_type":               jiraIssueTypeDataSourceType{},
		"atlassian_jira_issue_type_scheme":        jiraIssueTypeSchemeDataSourceType{},
		"atlassian_jira_issue_type_screen_scheme": &jiraIssueTypeScreenSchemeDataSourceType{},
		"atlassian_jira_screen_scheme":            &jiraScreenSchemeDataSourceType{},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
