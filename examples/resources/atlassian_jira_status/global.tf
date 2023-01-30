resource "atlassian_jira_status" "example" {
  name            = "foo"
  status_category = "TODO"
  status_scope = {
    type = "GLOBAL"
  }
}
