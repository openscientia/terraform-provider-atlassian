resource "atlassian_jira_issue_type_screen_scheme" "example" {
  name = "foo"
  issue_type_mappings = [
    {
      issue_type_id    = "default"
      screen_scheme_id = "10010"
    },
    {
      issue_type_id    = "10000" # id of epic issue type
      screen_scheme_id = "10100"
    },
    {
      issue_type_id    = "10001" # id of story issue type
      screen_scheme_id = "11000"
    }
  ]
}
