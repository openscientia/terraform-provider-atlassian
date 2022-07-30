package atlassian

import (
	"fmt"
	"strings"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraScreenScheme_Basic(t *testing.T) {
	resourceName := "atlassian_jira_screen_scheme.test"
	resourceAttributeName := sdkacctest.RandomWithPrefix("tf-test-screen-scheme")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScreenSchemeConfig_basic(resourceName, resourceAttributeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", resourceAttributeName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "screens.edit", "0"),
					resource.TestCheckResourceAttr(resourceName, "screens.create", "0"),
					resource.TestCheckResourceAttr(resourceName, "screens.view", "0"),
					resource.TestCheckResourceAttr(resourceName, "screens.default", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccJiraScreenScheme_Description(t *testing.T) {
	resourceName := "atlassian_jira_screen_scheme.test"
	resourceAttributeName := sdkacctest.RandomWithPrefix("tf-test-screen-scheme")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScreenSchemeConfig_description(resourceName, resourceAttributeName, "description1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				Config: testAccScreenSchemeConfig_description(resourceName, resourceAttributeName, "description2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}

func TestAccJiraScreenScheme_Screens(t *testing.T) {
	resourceName := "atlassian_jira_screen_scheme.test"
	resourceAttributeName := sdkacctest.RandomWithPrefix("tf-test-screen-scheme")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScreenSchemeConfig_screens(resourceName, resourceAttributeName, "1", "1", "1", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "screens.edit", "1"),
					resource.TestCheckResourceAttr(resourceName, "screens.create", "1"),
					resource.TestCheckResourceAttr(resourceName, "screens.view", "1"),
					resource.TestCheckResourceAttr(resourceName, "screens.default", "1"),
				),
			},
			{
				Config: testAccScreenSchemeConfig_screens(resourceName, resourceAttributeName, "2", "1", "1", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "screens.edit", "2"),
					resource.TestCheckResourceAttr(resourceName, "screens.create", "1"),
					resource.TestCheckResourceAttr(resourceName, "screens.view", "1"),
					resource.TestCheckResourceAttr(resourceName, "screens.default", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccScreenSchemeConfig_basic(resource_name, name string) string {
	splits := strings.Split(resource_name, ".")
	return fmt.Sprintf(
		`resource %[1]q %[2]q {
			name = %[3]q
			screens = {
				default = 1
			}
		}`, splits[0], splits[1], name,
	)
}

func testAccScreenSchemeConfig_description(resource_name, name, description string) string {
	splits := strings.Split(resource_name, ".")
	return fmt.Sprintf(
		`resource %[1]q %[2]q {
			name = %[3]q
			description = %[4]q
			screens = {
				default = 1
			}
		}`, splits[0], splits[1], name, description,
	)
}

func testAccScreenSchemeConfig_screens(resource_name, name string, screens ...string) string {
	splits := strings.Split(resource_name, ".")
	return fmt.Sprintf(
		`resource %[1]q %[2]q {
			name = %[3]q
			description = "description"
			screens = {
				edit = %[4]q
				create = %[5]q
				view = %[6]q
				default = %[7]q
			}
		}`, splits[0], splits[1], name, screens[0], screens[1], screens[2], screens[3],
	)
}
