resource "atlassian_jira_issue_field_configuration" "example" {
  name = "foo"
}

resource "atlassian_jira_issue_field_configuration_item" "example_text_renderer" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id       = "customfield_10000"
    renderer = "text-renderer"
  }
}

resource "atlassian_jira_issue_field_configuration_item" "example_wiki_renderer" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id       = "customfield_10001"
    renderer = "wiki-renderer"
  }
}
