package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJiraGroupUser_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-group-user")
	resourceName = "atlassian_jira_group_user.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupUserConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "group_name", "atlassian_jira_group.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "account_id", "data.atlassian_jira_myself.test", "account_id"),
					resource.TestCheckResourceAttrPair(resourceName, "self", "data.atlassian_jira_myself.test", "self"),
					resource.TestCheckResourceAttrPair(resourceName, "email_address", "data.atlassian_jira_myself.test", "email_address"),
					resource.TestCheckResourceAttrPair(resourceName, "display_name", "data.atlassian_jira_myself.test", "display_name"),
					resource.TestCheckResourceAttrPair(resourceName, "active", "data.atlassian_jira_myself.test", "active"),
					resource.TestCheckResourceAttrPair(resourceName, "active", "data.atlassian_jira_myself.test", "active"),
					resource.TestCheckResourceAttrPair(resourceName, "timezone", "data.atlassian_jira_myself.test", "timezone"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupUserImportConfig,
			},
		},
	})
}

func TestAccJiraGroupUser_ForceNewResource(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-group-user")
	resourceName = "atlassian_jira_group_user.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupUserConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", randomName),
				),
			},
			{
				Config: testAccGroupUserConfig_basic(resourceName, randomName+"b"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", randomName+"b"),
				),
			},
		},
	})
}

func testAccGroupUserImportConfig(s *terraform.State) (string, error) {
	group_name := s.RootModule().Resources["atlassian_jira_group_user.test"].Primary.Attributes["group_name"]
	account_id := s.RootModule().Resources["atlassian_jira_group_user.test"].Primary.Attributes["account_id"]
	return fmt.Sprintf("%s,%s", group_name, account_id), nil
}

func testAccGroupUserConfig_basic(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	data "atlassian_jira_myself" "test" {}

	resource "atlassian_jira_group" "test" { 
		name = %[3]q
	}

	resource %[1]q %[2]q {
		group_name = atlassian_jira_group.test.name
		account_id = data.atlassian_jira_myself.test.account_id
	}
	`, splits[0], splits[1], name)
}
