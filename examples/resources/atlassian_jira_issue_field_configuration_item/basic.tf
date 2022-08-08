resource "atlassian_jira_issue_field_configuration" "example" {
  name = "foo"
}

resource "atlassian_jira_issue_field_configuration_item" "example" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id = "customfield_10000"
  }
}
