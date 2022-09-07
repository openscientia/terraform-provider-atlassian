resource "atlassian_jira_permission_grant" "example" {
  permission_scheme_id = "10000"
  holder = {
    type      = "userCustomField"
    parameter = "10101"
  }
  permission = "ADD_COMMENTS"
}
