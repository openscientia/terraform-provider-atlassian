resource "atlassian_jira_issue_type_screen_scheme" "example" {
  name = "foo"
  issue_type_mappings = [
    {
      issue_type_id    = "default"
      screen_scheme_id = "10101"
    }
  ]
}
