package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraScreenSchemeDataSource_Basic(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("tf-test-screen-scheme")
	dataSourceName := "data.atlassian_jira_screen_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScreenSchemeDataSourceConfig_basic(dataSourceName, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", resourceName),
					resource.TestCheckResourceAttr(dataSourceName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceName, "screens.edit", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "screens.create", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "screens.view", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "screens.default", "1"),
				),
			},
		},
	})
}

func testAccScreenSchemeDataSourceConfig_basic(dataSourceName, resourceName string) string {
	splits := strings.Split(dataSourceName, ".")
	return fmt.Sprintf(`
  resource %[1]q %[2]q {
	name = %[3]q
	screens = {
		default = 1
	}
  }
  
  data %[1]q %[2]q {
	id = %[1]s.%[2]s.id
  }
  `, splits[1], splits[2], resourceName)
}
