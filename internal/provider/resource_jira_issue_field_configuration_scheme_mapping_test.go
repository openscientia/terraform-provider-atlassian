package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJiraIssueFieldConfigurationSchemeMapping_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-scheme-mapping")
	resourceName = "atlassian_jira_issue_field_configuration_scheme_mapping.test"
	issue_type_id := "10000" // epic
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationSchemeMapping_basic(resourceName, randomName, issue_type_id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_id", issue_type_id),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIssueFieldConfigurationSchemeMappingImportConfig,
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationSchemeMapping_FieldConfigurationSchemeId(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-scheme-mapping")
	resourceName = "atlassian_jira_issue_field_configuration_scheme_mapping.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationSchemeMapping_fieldconfigurationschemeid(resourceName, randomName, randomName+"A"),
			},
			{
				Config: testAccIssueFieldConfigurationSchemeMapping_fieldconfigurationschemeid(resourceName, randomName, randomName+"B"),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationSchemeMapping_FieldConfigurationId(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-scheme-mapping")
	resourceName = "atlassian_jira_issue_field_configuration_scheme_mapping.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationSchemeMapping_fieldconfigurationid(resourceName, randomName, randomName+"A"),
			},
			{
				Config: testAccIssueFieldConfigurationSchemeMapping_fieldconfigurationid(resourceName, randomName, randomName+"B"),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationSchemeMapping_IssueTypeId(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-scheme-mapping")
	resourceName = "atlassian_jira_issue_field_configuration_scheme_mapping.test"
	issue_type_id := []string{"10000", "10012"} // epic, story
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationSchemeMapping_issuetypeid(resourceName, randomName, issue_type_id[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_id", issue_type_id[0]),
				),
			},
			{
				Config: testAccIssueFieldConfigurationSchemeMapping_issuetypeid(resourceName, randomName, issue_type_id[1]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "issue_type_id", issue_type_id[1]),
				),
			},
		},
	})
}

func testAccIssueFieldConfigurationSchemeMapping_basic(resourceName, name, issueTypeId string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource "atlassian_jira_issue_field_configuration_scheme" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		field_configuration_scheme_id = atlassian_jira_issue_field_configuration_scheme.test.id
		field_configuration_id = atlassian_jira_issue_field_configuration.test.id
		issue_type_id = %[4]q
	}
	`, splits[0], splits[1], name, issueTypeId)
}

func testAccIssueFieldConfigurationSchemeMapping_fieldconfigurationschemeid(resourceName, name, fieldConfigurationSchemeName string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource "atlassian_jira_issue_field_configuration_scheme" %[4]q {
		name = %[4]q
	}

	resource %[1]q %[2]q {
		field_configuration_scheme_id = atlassian_jira_issue_field_configuration_scheme.%[4]s.id
		field_configuration_id = atlassian_jira_issue_field_configuration.test.id
		issue_type_id = "10000"
	}
	`, splits[0], splits[1], name, fieldConfigurationSchemeName)
}

func testAccIssueFieldConfigurationSchemeMapping_fieldconfigurationid(resourceName, name, fieldConfigurationName string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" %[4]q {
		name = %[4]q
	}

	resource "atlassian_jira_issue_field_configuration_scheme" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		field_configuration_scheme_id = atlassian_jira_issue_field_configuration_scheme.test.id
		field_configuration_id = atlassian_jira_issue_field_configuration.%[4]s.id
		issue_type_id = "10000"
	}
	`, splits[0], splits[1], name, fieldConfigurationName)
}

func testAccIssueFieldConfigurationSchemeMapping_issuetypeid(resourceName, name, issueTypeId string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource "atlassian_jira_issue_field_configuration_scheme" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		field_configuration_scheme_id = atlassian_jira_issue_field_configuration_scheme.test.id
		field_configuration_id = atlassian_jira_issue_field_configuration.test.id
		issue_type_id = %[4]q
	}
	`, splits[0], splits[1], name, issueTypeId)
}

func testAccIssueFieldConfigurationSchemeMappingImportConfig(s *terraform.State) (string, error) {
	fieldConfigurationSchemeId := s.RootModule().Resources["atlassian_jira_issue_field_configuration_scheme.test"].Primary.Attributes["id"]
	fieldConfigurationId := s.RootModule().Resources["atlassian_jira_issue_field_configuration.test"].Primary.Attributes["id"]
	issueTypeId := "10000" // epic
	return fmt.Sprintf("%s,%s,%s", fieldConfigurationSchemeId, fieldConfigurationId, issueTypeId), nil
}
