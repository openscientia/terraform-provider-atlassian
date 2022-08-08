package atlassian

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	resourceName = "atlassian_jira_issue_field_configuration_item.test"
)

func TestAccJiraIssueFieldConfigurationItem_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationItemConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "item.description"),
					resource.TestCheckResourceAttrSet(resourceName, "item.is_hidden"),
					resource.TestCheckResourceAttrSet(resourceName, "item.is_required"),
					resource.TestCheckResourceAttr(resourceName, "item.id", "customfield_10009"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIssueFieldConfigurationItemImportConfig,
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_Description(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationItemConfig_description(resourceName, randomName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.description", "description1"),
				),
			},
			{
				Config: testAccIssueFieldConfigurationItemConfig_description(resourceName, randomName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.description", "description2"),
				),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_IsHidden(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationItemConfig_ishidden(resourceName, randomName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.is_hidden", "false"),
				),
			},
			{
				Config: testAccIssueFieldConfigurationItemConfig_ishidden(resourceName, randomName, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.is_hidden", "true"),
				),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_IsRequired(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationItemConfig_isrequired(resourceName, randomName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.is_required", "false"),
				),
			},
			{
				Config: testAccIssueFieldConfigurationItemConfig_isrequired(resourceName, randomName, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.is_required", "true"),
				),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_Renderer(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationItemConfig_renderer(resourceName, randomName, "text-renderer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.renderer", "text-renderer"),
				),
			},
			{
				Config: testAccIssueFieldConfigurationItemConfig_renderer(resourceName, randomName, "wiki-renderer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.renderer", "wiki-renderer"),
				),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_RendererErrors(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccIssueFieldConfigurationItemConfig_renderererrors(resourceName, randomName, "autocomplete-renderer"),
				ExpectError: regexp.MustCompile(`Value must be one of`),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_NonRenderableErrors(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Actual End field
			{
				Config:      testAccIssueFieldConfigurationItemConfig_nonrenderableerrors(resourceName, randomName, "customfield_10009"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the non-renderable item"),
			},
			// Assignee field
			{
				Config:      testAccIssueFieldConfigurationItemConfig_nonrenderableerrors(resourceName, randomName, "assignee"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the non-renderable item"),
			},
			// Reporter field
			{
				Config:      testAccIssueFieldConfigurationItemConfig_nonrenderableerrors(resourceName, randomName, "reporter"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the non-renderable item"),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_LockedErrors(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Request Type
			{
				Config:      testAccIssueFieldConfigurationItemConfig_lockederrors(resourceName, randomName, "customfield_10010"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the locked item"),
			},
			// Epic Name field
			{
				Config:      testAccIssueFieldConfigurationItemConfig_lockederrors(resourceName, randomName, "customfield_10011"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the locked item"),
			},
			// Epic Color field
			{
				Config:      testAccIssueFieldConfigurationItemConfig_lockederrors(resourceName, randomName, "customfield_10013"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the locked item"),
			},
			// Epic Link field
			{
				Config:      testAccIssueFieldConfigurationItemConfig_nonrenderableerrors(resourceName, randomName, "customfield_10014"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the locked item"),
			},
			// Issue Color field
			{
				Config:      testAccIssueFieldConfigurationItemConfig_lockederrors(resourceName, randomName, "customfield_10017"),
				ExpectError: regexp.MustCompile("Tried to set a renderer for the locked item"),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_HiddenRequired(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationItemConfig_hiddenrequired(resourceName, randomName, "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.is_hidden", "true"),
					resource.TestCheckResourceAttr(resourceName, "item.is_required", "false"),
				),
			},
			{
				Config: testAccIssueFieldConfigurationItemConfig_hiddenrequired(resourceName, randomName, "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.is_hidden", "true"),
					resource.TestCheckResourceAttr(resourceName, "item.is_required", "true"),
				),
			},
		},
	})
}

func TestAccJiraIssueFieldConfigurationItem_ForceNewResource(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-field-configuration-item")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueFieldConfigurationItemConfig_forcenewresource(resourceName, randomName, "comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.id", "comment"),
				),
			},
			{
				Config: testAccIssueFieldConfigurationItemConfig_forcenewresource(resourceName, randomName+"2", "description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "item.id", "description"),
				),
			},
		},
	})
}

func testAccIssueFieldConfigurationItemImportConfig(s *terraform.State) (string, error) {
	issueFieldConfigurationID := s.RootModule().Resources["atlassian_jira_issue_field_configuration.test"].Primary.Attributes["id"]
	itemID := s.RootModule().Resources[resourceName].Primary.Attributes["item.id"]
	return fmt.Sprintf("%s,%s", issueFieldConfigurationID, itemID), nil
}

func testAccIssueFieldConfigurationItemConfig_basic(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = "customfield_10009"
		}
	}
	`, splits[0], splits[1], name)
}

func testAccIssueFieldConfigurationItemConfig_description(resourceName, name, description string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = "customfield_10009"
			description = %[4]q
		}
	}
	`, splits[0], splits[1], name, description)
}

func testAccIssueFieldConfigurationItemConfig_ishidden(resourceName, name, isHidden string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = "customfield_10009"
			is_hidden = %[4]s
		}
	}
	`, splits[0], splits[1], name, isHidden)
}

func testAccIssueFieldConfigurationItemConfig_isrequired(resourceName, name, isRequired string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = "customfield_10009"
			is_required = %[4]s
		}
	}
	`, splits[0], splits[1], name, isRequired)
}

func testAccIssueFieldConfigurationItemConfig_renderer(resourceName, name, renderer string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = "comment"
			renderer = %[4]q
		}
	}
	`, splits[0], splits[1], name, renderer)
}

func testAccIssueFieldConfigurationItemConfig_renderererrors(resourceName, name, renderer string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = "comment"
			renderer = %[4]q
		}
	}
	`, splits[0], splits[1], name, renderer)
}

func testAccIssueFieldConfigurationItemConfig_nonrenderableerrors(resourceName, name, itemID string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = %[4]q
			renderer = "text-renderer"
		}
	}
	`, splits[0], splits[1], name, itemID)
}

func testAccIssueFieldConfigurationItemConfig_lockederrors(resourceName, name, itemId string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = %[4]q
			renderer = "text-renderer"
		}
	}
	`, splits[0], splits[1], name, itemId)
}

func testAccIssueFieldConfigurationItemConfig_hiddenrequired(resourceName, name, isHidden, isRequired string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" "test" {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.test.id
		item = {
			id = "comment"
			is_hidden = %[4]s
			is_required = %[5]s
		}
	}
	`, splits[0], splits[1], name, isHidden, isRequired)
}

func testAccIssueFieldConfigurationItemConfig_forcenewresource(resourceName, name, itemId string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource "atlassian_jira_issue_field_configuration" %[3]q {
		name = %[3]q
	}

	resource %[1]q %[2]q {
		issue_field_configuration = atlassian_jira_issue_field_configuration.%[3]s.id
		item = {
			id = %[4]q
		}
	}
	`, splits[0], splits[1], name, itemId)
}
