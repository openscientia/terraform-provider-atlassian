package atlassian

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueScreenDataSource_Basic(t *testing.T) {
	dataSourceName := "data.atlassian_jira_issue_screen.test"
	r.Test(t, r.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []r.TestStep{
			// Read testing
			{
				Config: testAccJiraIssueScreenDataSourceConfig_basic,
				Check: r.ComposeAggregateTestCheckFunc(
					r.TestCheckResourceAttr(dataSourceName, "name", "Default Screen"),
				),
			},
		},
	})
}

func TestAccJiraIssueScreenDataSource_ErrorCases(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []r.TestStep{
			{
				Config: `
					data "atlassian_jira_issue_screen" "test" {
					}
				`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config: `
					data "atlassian_jira_issue_screen" "test" {
						id = "foo"
					}
				`,
				ExpectError: regexp.MustCompile("Unable to parse value of \"id\" attribute"),
			},
		},
	})
}

const testAccJiraIssueScreenDataSourceConfig_basic = `
data "atlassian_jira_issue_screen" "test" {
  id = "1" // id of default screen
}
`
