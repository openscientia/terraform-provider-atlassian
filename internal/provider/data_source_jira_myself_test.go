package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraMyselfDataSource_Basic(t *testing.T) {
	dataSourceName := "data.atlassian_jira_myself.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMyselfDataSourceConfig_basic(dataSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "account_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", dataSourceName, "account_id"),
				),
			},
		},
	})
}

func testAccMyselfDataSourceConfig_basic(dataSourceName string) string {
	splits := strings.Split(dataSourceName, ".")
	return fmt.Sprintf(`
	  data %[1]q %[2]q {}
	`, splits[1], splits[2])
}
