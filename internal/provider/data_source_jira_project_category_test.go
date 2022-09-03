package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraProjectCategoryDataSource_Basic(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("tf-test-project-category")
	dataSourceName := "data.atlassian_jira_project_category.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectCategoryDataSourceConfig_basic(dataSourceName, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", resourceName),
					resource.TestCheckResourceAttr(dataSourceName, "description", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "self"),
				),
			},
		},
	})
}

func testAccProjectCategoryDataSourceConfig_basic(dataSourceName, name string) string {
	splits := strings.Split(dataSourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		name = %[3]q
	  }
	  
	  data %[1]q %[2]q {
		id = %[1]s.%[2]s.id
	  }
	`, splits[1], splits[2], name)
}
