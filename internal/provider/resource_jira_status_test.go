package atlassian

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraStatus_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", " "),
					resource.TestCheckResourceAttr(resourceName, "status_category", "TODO"),
					resource.TestCheckResourceAttr(resourceName, "status_scope.type", "GLOBAL"),
					resource.TestCheckResourceAttr(resourceName, "status_scope.id", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

func TestAccJiraStatus_Name(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusConfig_basic(resourceName, randomName+"1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName+"1"),
				),
			},
			{
				Config: testAccStatusConfig_basic(resourceName, randomName+"2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName+"2"),
				),
			},
		},
	})
}

func TestAccJiraStatus_Description(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusConfig_description(resourceName, randomName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				Config: testAccStatusConfig_description(resourceName, randomName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func TestAccJiraStatus_StatusCategory(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusConfig_statuscategory(resourceName, randomName, "TODO"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status_category", "TODO"),
				),
			},
			{
				Config: testAccStatusConfig_statuscategory(resourceName, randomName, "DONE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status_category", "DONE"),
				),
			},
		},
	})
}

func TestAccJiraStatus_StatusScopeType(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status_scope.type", "GLOBAL"),
				),
			},
			{
				Config: testAccStatusConfig_statusscopetype(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status_scope.type", "PROJECT"),
				),
			},
		},
	})
}

func TestAccJiraStatus_StatusScopeTypeErrors(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccStatusConfig_statusScopeTypeErrors(resourceName, randomName),
				ExpectError: regexp.MustCompile(`must not have a value for "status_scope.id" attribute`),
			},
		},
	})
}

func TestAccJiraStatus_StatusScopeId(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusConfig_statusscopeid(resourceName, randomName, "10003"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status_scope.id", "10003"),
				),
			},
			{
				Config: testAccStatusConfig_statusscopeid(resourceName, randomName, "10004"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status_scope.id", "10004"),
				),
			},
		},
	})
}

func TestAccJiraStatus_StatusScopeIdErrors(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-jira-status")
	resourceName = "atlassian_jira_status.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccStatusConfig_statusScopeIdErrors(resourceName, randomName),
				ExpectError: regexp.MustCompile(`Failed to provide value for "status_scope.id" attribute`),
			},
		},
	})
}

func testAccStatusConfig_basic(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		status_category = "TODO"
		status_scope = {
			type = "GLOBAL"
		}
	}
	`, splits[0], splits[1], name)
}

func testAccStatusConfig_description(resourceName, name, description string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		description = %[4]q
		status_category = "TODO"
		status_scope = {
			type = "GLOBAL"
		}
	}
	`, splits[0], splits[1], name, description)
}

func testAccStatusConfig_statuscategory(resourceName, name, status_category string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		status_category = %[4]q
		status_scope = {
			type = "GLOBAL"
		}
	}
	`, splits[0], splits[1], name, status_category)
}

func testAccStatusConfig_statusscopetype(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		status_category = "TODO"
		status_scope = {
			type = "PROJECT"
			id = "10003"
		}
	}
	`, splits[0], splits[1], name)
}

func testAccStatusConfig_statusScopeTypeErrors(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		status_category = "TODO"
		status_scope = {
			type = "GLOBAL"
			id = "10001"
		}
	}
	`, splits[0], splits[1], name)
}

func testAccStatusConfig_statusscopeid(resourceName, name, statusScopeId string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		status_category = "TODO"
		status_scope = {
			type = "PROJECT"
			id = %[4]q
		}
	}
	`, splits[0], splits[1], name, statusScopeId)
}

func testAccStatusConfig_statusScopeIdErrors(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		status_category = "TODO"
		status_scope = {
			type = "PROJECT"
		}
	}
	`, splits[0], splits[1], name)
}
