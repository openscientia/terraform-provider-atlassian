resource "atlassian_jira_permission_grant" "example" {
  permission_scheme_id = "10000"
  holder = {
    type      = "projectRole"
    parameter = "10001"
  }
  permission = "ADMINISTER_PROJECTS"
}
