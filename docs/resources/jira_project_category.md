---
page_title: "Atlassian Cloud: atlassian_jira_project_category"
subcategory: "Jira Cloud"
description: |-
  Manages atlassian_jira_project_category.
---

# Resource: atlassian_jira_project_category

Provides an `atlassian_jira_project_category` resource.

Learn more about [Jira Project Categories](https://support.atlassian.com/jira-cloud-administration/docs/add-assign-and-delete-project-categories/).

See more details about the [Jira Cloud Platform REST API for Project Categories](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-categories/#api-group-project-categories).

## Example Usage

### Basic

```terraform
resource "atlassian_jira_project_category" "example" {
  name = "foo"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the project category. The name must be unique. The maximum length is 255 characters.

### Optional

- `description` (String) The description of the project category. The maximum length is 1000 characters.
- `self` (String) The URL of the project category.

### Read-Only

- `id` (String) The ID of the project category.

## Import

`atlassian_jira_project_category` can be imported using `id`, e.g.,

```sh
$ terraform import atlassian_jira_project_category.foo 1234567890
```
