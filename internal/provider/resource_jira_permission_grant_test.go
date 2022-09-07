package atlassian

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJiraPermissionGrant_Basic(t *testing.T) {
	resourceName = "atlassian_jira_permission_grant.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionGrantConfig_basic(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "permission_scheme_id", "10004"),
					resource.TestCheckResourceAttr(resourceName, "holder.type", "anyone"),
					resource.TestCheckResourceAttr(resourceName, "holder.parameter", ""),
					resource.TestCheckResourceAttr(resourceName, "permission", "ADMINISTER_PROJECTS"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPermissionGrantImportConfig,
			},
		},
	})
}

func TestAccJiraPermissionGrant_PermissionSchemeId(t *testing.T) {
	resourceName = "atlassian_jira_permission_grant.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionGrantConfig_permissionschemeid(resourceName, "10004"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "permission_scheme_id", "10004"),
				),
			},
			{
				Config: testAccPermissionGrantConfig_permissionschemeid(resourceName, "10005"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "permission_scheme_id", "10005"),
				),
			},
		},
	})
}

func TestAccJiraPermissionGrant_HolderType(t *testing.T) {
	resourceName = "atlassian_jira_permission_grant.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionGrantConfig_holdertype(resourceName, "projectLead"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "holder.type", "projectLead"),
				),
			},
			{
				Config: testAccPermissionGrantConfig_holdertype(resourceName, "reporter"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "holder.type", "reporter"),
				),
			},
		},
	})
}

func TestAccJiraPermissionGrant_HolderTypeErrors(t *testing.T) {
	resourceName = "atlassian_jira_permission_grant.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccPermissionGrantConfig_holdertypeerrors(resourceName, "group"),
				ExpectError: regexp.MustCompile(`Failed to provide a value`),
			},
			{
				Config:      testAccPermissionGrantConfig_holdertypeerrors(resourceName, "projectRole"),
				ExpectError: regexp.MustCompile(`Failed to provide a value`),
			},
			{
				Config:      testAccPermissionGrantConfig_holdertypeerrors(resourceName, "user"),
				ExpectError: regexp.MustCompile(`Failed to provide a value`),
			},
			{
				Config:      testAccPermissionGrantConfig_holdertypeerrors(resourceName, "userCustomField"),
				ExpectError: regexp.MustCompile(`Failed to provide a value`),
			},
		},
	})
}

func TestAccJiraPermissionGrant_HolderParameter(t *testing.T) {
	resourceName = "atlassian_jira_permission_grant.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionGrantConfig_holderparameter(resourceName, "administrators"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "holder.parameter", "administrators"),
				),
			},
			{
				Config: testAccPermissionGrantConfig_holderparameter(resourceName, "site-admins"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "holder.parameter", "site-admins"),
				),
			},
		},
	})
}

func TestAccJiraPermissionGrant_Permission(t *testing.T) {
	resourceName = "atlassian_jira_permission_grant.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionGrantConfig_permission(resourceName, "CREATE_ISSUES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "permission", "CREATE_ISSUES"),
				),
			},
			{
				Config: testAccPermissionGrantConfig_permission(resourceName, "DELETE_ISSUES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "permission", "DELETE_ISSUES"),
				),
			},
		},
	})
}

func testAccPermissionGrantImportConfig(s *terraform.State) (string, error) {
	permissionGrantId := s.RootModule().Resources["atlassian_jira_permission_grant.test"].Primary.Attributes["id"]
	permissionSchemeId := s.RootModule().Resources["atlassian_jira_permission_grant.test"].Primary.Attributes["permission_scheme_id"]
	return fmt.Sprintf("%s,%s", permissionGrantId, permissionSchemeId), nil
}

func testAccPermissionGrantConfig_basic(resourceName string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		permission_scheme_id = "10004"
		holder = {
			type = "anyone"
		}
		permission = "ADMINISTER_PROJECTS"
	}
	`, splits[0], splits[1])
}

func testAccPermissionGrantConfig_permissionschemeid(resourceName, schemeId string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		permission_scheme_id = %[3]q
		holder = {
			type = "assignee"
		}
		permission = "BROWSE_PROJECTS"
	}
	`, splits[0], splits[1], schemeId)
}

func testAccPermissionGrantConfig_holdertype(resourceName, holderType string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		permission_scheme_id = "10004"
		holder = {
			type = %[3]q
		}
		permission = "ASSIGN_ISSUES"
	}
	`, splits[0], splits[1], holderType)
}

func testAccPermissionGrantConfig_holderparameter(resourceName, holderParameter string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		permission_scheme_id = "10004"
		holder = {
			type = "group"
			parameter = %[3]q
		}
		permission = "CLOSE_ISSUES"
	}
	`, splits[0], splits[1], holderParameter)
}

func testAccPermissionGrantConfig_permission(resourceName, permission string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		permission_scheme_id = "10004"
		holder = {
			type = "projectLead"
		}
		permission = %[3]q
	}
	`, splits[0], splits[1], permission)
}

func testAccPermissionGrantConfig_holdertypeerrors(resourceName, holderType string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		permission_scheme_id = "10004"
		holder = {
			type = %[3]q
		}
		permission = "EDIT_ISSUES"
	}
	`, splits[0], splits[1], holderType)
}
