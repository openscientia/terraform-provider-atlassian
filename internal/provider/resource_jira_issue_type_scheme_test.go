package atlassian

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeSchemeResource(t *testing.T) {
	resourceName := "atlassian_jira_issue_type_scheme.test"
	testAttributeNames := []string{"Test Issue Type Scheme 1", "Test Jira Issue Type Scheme 2"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraIssueTypeSchemeResourceConfig(testAttributeNames[0]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", testAttributeNames[0]),
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
				Config: testAccJiraIssueTypeSchemeResourceConfig(testAttributeNames[1]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testAttributeNames[1]),
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
