package atlassian

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJiraIssueType_Basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type")
	resourceName := "atlassian_jira_issue_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "hierarchy_level", "0"),
					resource.TestCheckResourceAttr(resourceName, "avatar_id", "10300"),
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

func TestAccJiraIssueType_Name(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type")
	resourceName := "atlassian_jira_issue_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeConfig_basic(resourceName, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
				),
			},
			{
				Config: testAccJiraIssueTypeConfig_basic(resourceName, randomName+"B"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName+"B"),
				),
			},
		},
	})
}

func TestAccJiraIssueType_Description(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type")
	resourceName := "atlassian_jira_issue_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeConfig_description(resourceName, randomName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				Config: testAccJiraIssueTypeConfig_description(resourceName, randomName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func TestAccJiraIssueType_Type(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type")
	resourceName := "atlassian_jira_issue_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeConfig_type(resourceName, randomName, "standard"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "hierarchy_level", "0"),
				),
			},
			{
				Config: testAccJiraIssueTypeConfig_type(resourceName+"b", randomName+"b", "sub-task"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName+"b", "type", "sub-task"),
					resource.TestCheckResourceAttr(resourceName+"b", "hierarchy_level", "-1"),
				),
			},
		},
	})
}

func TestAccJiraIssueType_HierarchyLevel(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type")
	resourceName := "atlassian_jira_issue_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeConfig_hierarchylevel(resourceName, randomName, "0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hierarchy_level", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "standard"),
				),
			},
			{
				Config: testAccJiraIssueTypeConfig_hierarchylevel(resourceName+"b", randomName+"b", "-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName+"b", "hierarchy_level", "-1"),
					resource.TestCheckResourceAttr(resourceName+"b", "type", "sub-task"),
				),
			},
		},
	})
}

func TestAccJiraIssueType_AvatarId(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-issue-type")
	resourceName := "atlassian_jira_issue_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueTypeConfig_avatarid(resourceName, randomName, "10300"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "avatar_id", "10300"),
				),
			},
			{
				Config: testAccJiraIssueTypeConfig_avatarid(resourceName, randomName, "10315"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "avatar_id", "10315"),
				),
			},
		},
	})
}

func testAccJiraIssueTypeConfig_basic(resourceName, name string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
  		name = %[3]q
	}`, splits[0], splits[1], name)
}

func testAccJiraIssueTypeConfig_description(resourceName, name, description string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
  		name = %[3]q
		description = %[4]q
	}`, splits[0], splits[1], name, description)
}

func testAccJiraIssueTypeConfig_type(resourceName, name, types string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
  		name = %[3]q
		type = %[4]q
	}`, splits[0], splits[1], name, types)
}

func testAccJiraIssueTypeConfig_hierarchylevel(resourceName, name, hierarchy_level string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
  		name = %[3]q
		hierarchy_level = %[4]q
	}`, splits[0], splits[1], name, hierarchy_level)
}

func testAccJiraIssueTypeConfig_avatarid(resourceName, name, avatar_id string) string {
	splits := strings.Split(resourceName, ".")
	return fmt.Sprintf(`
	resource %[1]q %[2]q {
  		name = %[3]q
		avatar_id = %[4]q
	}`, splits[0], splits[1], name, avatar_id)
}
