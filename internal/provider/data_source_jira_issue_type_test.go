package atlassian

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeDataSource(t *testing.T) {
	dataSourceName := "data.atlassian_jira_issue_type.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccJiraIssueTypeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "Epic"),
				),
			},
		},
	})
}

const testAccJiraIssueTypeDataSourceConfig = `
data "atlassian_jira_issue_type" "test" {
  id = "10000" // default id of epic issue type
}
`
