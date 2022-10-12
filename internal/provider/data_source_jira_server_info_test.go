package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraServerInfoDataSource_Basic(t *testing.T) {
	dataSourceName := "data.atlassian_jira_server_info.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServerInfoDataSourceConfig_basic(dataSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "base_url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_numbers.0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_numbers.1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_numbers.2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployment_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "build_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "build_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scm_info"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_title"),
				),
			},
		},
	})
}

func testAccServerInfoDataSourceConfig_basic(dataSourceName string) string {
	splits := strings.Split(dataSourceName, ".")
	return fmt.Sprintf(`
	  data %[1]q %[2]q {}
	`, splits[1], splits[2])
}
