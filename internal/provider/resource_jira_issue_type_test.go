package atlassian

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeResource(t *testing.T) {
	resourceName := "atlassian_jira_issue_type.test"
	testAtrributeNames := []string{"Test Jira Issue Type 1", "Test Jira Issue Type 2"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraIssueTypeResourceConfig(testAtrributeNames[0]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", testAtrributeNames[0]),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "hierarchy_level", "0"),
					resource.TestCheckResourceAttr(resourceName, "avatar_id", "10300"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccJiraIssueTypeResourceConfig(testAtrributeNames[1]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testAtrributeNames[1]),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccJiraIssueTypeResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "atlassian_jira_issue_type" "test" {
  name = %[1]q
}
`, name)
}
