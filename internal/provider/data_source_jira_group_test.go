package atlassian

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraGroupDataSource_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-group")
	dataSourceName := "data.atlassian_jira_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_basic(dataSourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("atlassian_jira_group.test", "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair("atlassian_jira_group.test", "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair("atlassian_jira_group.test", "group_id", dataSourceName, "group_id"),
					resource.TestCheckResourceAttrPair("atlassian_jira_group.test", "self", dataSourceName, "self"),
					resource.TestCheckResourceAttrPair("atlassian_jira_group.test", "users.#", dataSourceName, "users.#"),
				),
			},
		},
	})
}

func TestAccJiraGroupDataSource_Users(t *testing.T) {
	dataSourceName := "data.atlassian_jira_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_users(dataSourceName, "site-admins"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceName, "users.*", map[string]*regexp.Regexp{
						"self":          regexp.MustCompile(`^https://[a-zA-Z0-9_\-\.]+\.atlassian\.net/rest/api/3/user\?accountId=[a-z0-9]{0,128}$`),
						"account_id":    regexp.MustCompile(`[a-z0-9]{0,128}`),
						"email_address": regexp.MustCompile(`.*`),
						"display_name":  regexp.MustCompile(`.*`),
						"active":        regexp.MustCompile(`true|false`),
						"timezone":      regexp.MustCompile(`[a-zA-Z0-9/]+|\z`),
						"account_type":  regexp.MustCompile(`atlassian|app|customer`),
					}),
				),
			},
		},
	})
}

func testAccGroupDataSourceConfig_basic(dataSourceName, name string) string {
	splits := strings.Split(dataSourceName, ".")
	return fmt.Sprintf(`
	  resource %[1]q %[2]q {
		name = %[3]q
	  }

	  data %[1]q %[2]q {
		name = %[1]s.%[2]s.name
	  }
	`, splits[1], splits[2], name)
}

func testAccGroupDataSourceConfig_users(dataSourceName, name string) string {
	splits := strings.Split(dataSourceName, ".")
	return fmt.Sprintf(`
	  data %[1]q %[2]q {
		name = %[3]q
	  }
	`, splits[1], splits[2], name)
}
