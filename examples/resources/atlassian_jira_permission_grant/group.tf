resource "atlassian_jira_permission_grant" "example" {
  permission_scheme_id = "10000"
  holder = {
    type      = "group"
    parameter = "site-admins"
  }
  permission = "BROWSE_PROJECTS"
}
