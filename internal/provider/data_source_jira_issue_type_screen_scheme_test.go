package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueTypeScreenSchemeDataSource_Basic(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("tf-test-issue-type-screen-scheme")
	dataSourceName := "data.atlassian_jira_issue_type_screen_scheme.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueTypeScreenSchemeDataSourceConfig_basic(dataSourceName, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", resourceName),
					resource.TestCheckResourceAttr(dataSourceName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceName, "issue_type_mappings.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "issue_type_mappings.0.issue_type_id", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "issue_type_mappings.0.screen_scheme_id", "1"),
				),
			},
		},
	})
}

func testAccIssueTypeScreenSchemeDataSourceConfig_basic(dataSourceName, name string) string {
	splits := strings.Split(dataSourceName, ".")
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
	  
	  data %[1]q %[2]q {
		id = %[1]s.%[2]s.id
	  }
	`, splits[1], splits[2], name)
}
