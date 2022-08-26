package atlassian

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeResource(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type")
	resourceName := "atlassian_jira_issue_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraIssueTypeResourceConfig(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
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
				Config: testAccJiraIssueTypeResourceConfig(randomName + "B"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName+"B"),
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
