package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeScheme_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-scheme")
	resourceName := "atlassian_jira_issue_type_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeSchemeConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "issue_type_ids.#", "1"),
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

func TestAccJiraIssueTypeScheme_Description(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-scheme")
	resourceName := "atlassian_jira_issue_type_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeSchemeConfig_description(resourceName, randomName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				Config: testAccJiraIssueTypeSchemeConfig_description(resourceName, randomName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func TestAccJiraIssueTypeScheme_IssueTypeIds(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type-scheme")
	resourceName := "atlassian_jira_issue_type_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeSchemeConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_ids.#", "1"),
				),
			},
			{
				Config: testAccJiraIssueTypeSchemeConfig_issuetypeids(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_ids.#", "2"),
				),
			},
		},
	})
}

func testAccJiraIssueTypeSchemeConfig_basic(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_type" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		name = %[3]q
		issue_type_ids = [resource.atlassian_jira_issue_type.test.id]
	}
	`, splits[0], splits[1], name)
}

func testAccJiraIssueTypeSchemeConfig_description(resourceName, name, description string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_type" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		name = %[3]q
		description = %[4]q
		issue_type_ids = [resource.atlassian_jira_issue_type.test.id]
	}
	`, splits[0], splits[1], name, description)
}

func testAccJiraIssueTypeSchemeConfig_issuetypeids(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_type" "test" {
		name = %[3]q
	}

	resource "atlassian_jira_issue_type" "testb" {
		name = "%[3]s2"
	}

	resource %[1]q %[2]q {
		name = %[3]q
		issue_type_ids = [resource.atlassian_jira_issue_type.test.id, resource.atlassian_jira_issue_type.testb.id]
	}
	`, splits[0], splits[1], name)
}
