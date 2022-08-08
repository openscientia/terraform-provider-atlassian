resource "atlassian_jira_issue_field_configuration" "example" {
  name = "foo"
}

resource "atlassian_jira_issue_field_configuration_item" "example_hide" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id        = "customfield_10000"
    is_hidden = true
  }
}

resource "atlassian_jira_issue_field_configuration_item" "example_show" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id        = "customfield_10001"
    is_hidden = false
  }
}
