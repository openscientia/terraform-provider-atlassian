package atlassian

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeSchemeDataSource(t *testing.T) {
	dataSourceName := "data.atlassian_jira_issue_type_scheme.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccJiraIssueTypeSchemeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "Default Issue Type Scheme"),
				),
			},
		},
	})
}

const testAccJiraIssueTypeSchemeDataSourceConfig = `
data "atlassian_jira_issue_type_scheme" "test" {
  id = "10000" // id of default issue type scheme
}
`
