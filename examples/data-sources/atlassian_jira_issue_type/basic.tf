data "atlassian_jira_issue_type" "example" {
  id = "10000" // default id of epic issue type
}

resource "atlassian_jira_issue_type_scheme" "example" {
  name           = "Example Issue Type Scheme"
  issue_type_ids = [data.atlassian_jira_issue_type.example.id]
}
