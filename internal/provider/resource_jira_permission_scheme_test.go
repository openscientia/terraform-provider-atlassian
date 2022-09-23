package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraPermissionScheme_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-permission-scheme")
	resourceName := "atlassian_jira_permission_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionScheme_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "self"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func TestAccJiraPermissionScheme_Name(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-permission-scheme")
	resourceName := "atlassian_jira_permission_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionScheme_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
				),
			},
			{
				Config: testAccPermissionScheme_basic(resourceName, randomName+"2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName+"2"),
				),
			},
		},
	})
}

func TestAccJiraPermissionScheme_Description(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-permission-scheme")
	resourceName := "atlassian_jira_permission_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionScheme_description(resourceName, randomName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				Config: testAccPermissionScheme_description(resourceName, randomName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func testAccPermissionScheme_basic(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
	}
	`, splits[0], splits[1], name)
}

func testAccPermissionScheme_description(resourceName, name, description string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
		description = %[4]q
	}
	`, splits[0], splits[1], name, description)
}
