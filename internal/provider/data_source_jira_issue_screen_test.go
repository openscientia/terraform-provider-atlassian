package atlassian

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueScreen_DataSource_Basic(t *testing.T) {
	dataSourceName := "data.atlassian_jira_issue_screen.test"
	r.Test(t, r.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []r.TestStep{
			// Read testing
			{
				Config: testAccJiraIssueScreenDataSourceConfig,
				Check: r.ComposeAggregateTestCheckFunc(
					r.TestCheckResourceAttr(dataSourceName, "name", "Default Screen"),
				),
			},
		},
	})
}

const testAccJiraIssueScreenDataSourceConfig = `
data "atlassian_jira_issue_screen" "test" {
  id = "1" // id of default screen
}
`
