package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraPermissionGrantDataSource_Basic(t *testing.T) {
	dataSourceName := "data.atlassian_jira_permission_grant.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionGrantDataSourceConfig_basic(dataSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "permission_scheme_id", "10004"),
					resource.TestCheckResourceAttr(dataSourceName, "holder.type", "assignee"),
					resource.TestCheckResourceAttr(dataSourceName, "permission", "CREATE_ISSUES"),
				),
			},
		},
	})
}

func testAccPermissionGrantDataSourceConfig_basic(dataSourceName string) string {
	splits := strings.Split(dataSourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
		permission_scheme_id = "10004"
		holder = {
			type = "assignee"
		}
		permission = "CREATE_ISSUES"
	}
	data %[1]q %[2]q {
		id = %[1]s.%[2]s.id
		permission_scheme_id = "10004"
	}
	`, splits[1], splits[2])
}
