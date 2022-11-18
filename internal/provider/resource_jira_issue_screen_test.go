package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueScreen_Basic(t *testing.T) {
	resourceName := "atlassian_jira_issue_screen.test"
	testAttributeNames := []string{"Test Issue Screen 1", "Test Issue Screen 2"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueScreenConfig_basic(resourceName, testAttributeNames[0]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", testAttributeNames[0]),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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
				Config: testAccJiraIssueScreenConfig_basic(resourceName, testAttributeNames[1]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testAttributeNames[1]),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccJiraIssueScreenConfig_basic(resource_name, name string) string {
	splits := strings.Split(resource_name, ".")
	return fmt.Sprintf(
		`resource %[1]q %[2]q {
			name = %[3]q
		}`, splits[0], splits[1], name,
	)
}
