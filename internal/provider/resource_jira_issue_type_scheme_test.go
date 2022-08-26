package atlassian

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeSchemeResource(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-scheme")
	resourceName := "atlassian_jira_issue_type_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraIssueTypeSchemeResourceConfig(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "issue_type_ids.#", "1"),
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
				Config: testAccJiraIssueTypeSchemeResourceConfig(randomName + "B"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName+"B"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccJiraIssueTypeSchemeResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "atlassian_jira_issue_type" "foo" {
  name = "[TF] Jira Issue Type Foo"
}

resource "atlassian_jira_issue_type_scheme" "test" {
	name = %[1]q
	issue_type_ids = [resource.atlassian_jira_issue_type.foo.id]
}
`, name)
}
