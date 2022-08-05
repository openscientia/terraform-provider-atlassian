package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeScreenScheme_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-screen-scheme")
	resourceName := "atlassian_jira_issue_type_screen_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueTypeScreenSchemeConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.issue_type_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.screen_scheme_id", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccJiraIssueTypeScreenScheme_Description(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-screen-scheme")
	resourceName := "atlassian_jira_issue_type_screen_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueTypeScreenSchemeConfig_description(resourceName, randomName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				Config: testAccIssueTypeScreenSchemeConfig_description(resourceName, randomName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func TestAccJiraIssueTypeScreenScheme_DefaultMapping(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-screen-scheme")
	resourceName := "atlassian_jira_issue_type_screen_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueTypeScreenSchemeConfig_defaultmapping(resourceName, randomName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.issue_type_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.screen_scheme_id", "1"),
				),
			},
			{
				Config: testAccIssueTypeScreenSchemeConfig_defaultmapping(resourceName, randomName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.issue_type_id", "default"),
				),
			},
		},
	})
}

func TestAccJiraIssueTypeScreenScheme_IssueTypeMappings(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-screen-scheme")
	resourceName := "atlassian_jira_issue_type_screen_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueTypeScreenSchemeConfig_issuetypemappings(resourceName, randomName, "10000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.issue_type_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.screen_scheme_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.1.issue_type_id", "10000"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.1.screen_scheme_id", "1"),
				),
			},
			{
				Config: testAccIssueTypeScreenSchemeConfig_issuetypemappings(resourceName, randomName, "10001"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.issue_type_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.0.screen_scheme_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.1.issue_type_id", "10001"),
					resource.TestCheckResourceAttr(resourceName, "issue_type_mappings.1.screen_scheme_id", "1"),
				),
			},
		},
	},
	)
}

func testAccIssueTypeScreenSchemeConfig_basic(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		issue_type_mappings = [
			{
				issue_type_id = "default"
				screen_scheme_id = "1"
			}
		]
	}
	`, splits[0], splits[1], name)
}

func testAccIssueTypeScreenSchemeConfig_description(resourceName, name, description string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		description = %[4]q
		issue_type_mappings = [
			{
				issue_type_id = "default"
				screen_scheme_id = "1"
			}
		]
	}
	`, splits[0], splits[1], name, description)
}

func testAccIssueTypeScreenSchemeConfig_defaultmapping(resourceName, name string, useDefaultScreenScheme bool) string {
	splits := strings.Split(resourceName, ".")
	if useDefaultScreenScheme {
		return fmt.Sprintf(`
		resource %[1]q %[2]q {
			name = %[3]q
			issue_type_mappings = [
				{
					issue_type_id = "default"
					screen_scheme_id = "1"
				}
			]
		}
		`, splits[0], splits[1], name)
	} else {
		return fmt.Sprintf(`
		resource "atlassian_jira_screen_scheme" "test" {
			name = %[3]q
			screens = {
				default = 1
			}
		}

		resource %[1]q %[2]q {
			name = %[3]q
			issue_type_mappings = [
				{
					issue_type_id = "default"
					screen_scheme_id = atlassian_jira_screen_scheme.test.id
				}
			]
		}
		`, splits[0], splits[1], name)
	}
}

func testAccIssueTypeScreenSchemeConfig_issuetypemappings(resourceName, name, nonDefaultIssueTypeId string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		issue_type_mappings = [
			{
				issue_type_id = "default"
				screen_scheme_id = "1"
			},
			{
				issue_type_id = %[4]q 
				screen_scheme_id = "1" 
			}
		]
	}
	`, splits[0], splits[1], name, nonDefaultIssueTypeId)
}
