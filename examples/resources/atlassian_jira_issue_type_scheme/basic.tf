resource "atlassian_jira_issue_type" "example" {
  name = "bar"
}

resource "atlassian_jira_issue_type_scheme" "example" {
  name           = "Example Jira Issue Type Scheme"
  issue_type_ids = [resource.atlassian_jira_issue_type.example.id]
}
