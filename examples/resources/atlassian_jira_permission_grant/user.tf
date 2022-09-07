resource "atlassian_jira_permission_grant" "example" {
  permission_scheme_id = "10000"
  holder = {
    type      = "user"
    parameter = "09876543a21b0c1234567d89"
  }
  permission = "ASSIGN_ISSUES"
}
