resource "atlassian_jira_screen_scheme" "example" {
  name = "foo"
  screens = {
    default = 1 # id of default screen scheme
  }
}
