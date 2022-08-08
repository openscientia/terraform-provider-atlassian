---
page_title: "Atlassian Cloud: atlassian_jira_issue_field_configuration_item"
subcategory: "Jira Cloud"
description: |-
  Manages atlassian_jira_issue_field_configuration_item.
---

# Resource: atlassian_jira_issue_field_configuration_item

Provides an `atlassian_jira_issue_field_configuration_item` resource.

Learn more about [Jira Issue Field Configuration Items](https://support.atlassian.com/jira-cloud-administration/docs/change-a-field-configuration/).

See more details about the [Jira Cloud Platform REST API for Issue Field Configuration Items](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-field-configurations/#api-rest-api-3-fieldconfiguration-id-fields-get).

~> **Note** `terraform destroy` does not delete `atlassian_jira_issue_field_configuration_item` but does remove the resource from Terraform state.

-> **Note** `atlassian_jira_issue_field_configuration_item` can only reference [`atlassian_jira_issue_field_configuration`](https://registry.terraform.io/providers/openscientia/atlassian/latest/docs/resources/jira_issue_field_configuration) used in [company-managed (classic) projects](https://support.atlassian.com/jira-software-cloud/docs/what-are-team-managed-and-company-managed-projects/).

## Example Usage

### Basic

```terraform
resource "atlassian_jira_issue_field_configuration" "example" {
  name = "foo"
}

resource "atlassian_jira_issue_field_configuration_item" "example" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id = "customfield_10000"
  }
}
```

### Hide or Show fields

-> **Note** `item.is_hidden` can be set to `true` to ensure that the field does not appear on any `atlassian_jira_issue_screen` (i.e. issue operation screens, workflow transition screens) where a specific `atlassian_jira_issue_field_configuration` applies. See more [details](https://support.atlassian.com/jira-cloud-administration/docs/specify-field-behavior/#Hide-or-show-a-field).

~> **Note** Hiding a field in the `atlassian_jira_issue_field_configuration_item` resource is distinct from not adding a field to a `atlassian_jira_issue_screen`. Fields hidden through the `atlassian_jira_issue_field_configuration_item` resource will be hidden in all applicable screens, regardless of whether or not they have been added to any `atlassian_jira_issue_screen`.

```terraform
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
```

### Required or Optional fields

-> **Note** `item.is_required` can be set to `true` to ensures that Jira validates the field has been given a value whenever an issue is edited. Alternatively, set to `false` to make a field optional if no value is required. See more [details](https://support.atlassian.com/jira-cloud-administration/docs/specify-field-behavior/#Make-a-field-required-or-optional).

~> **Note** If you make a field required, ensure that the field is present on all Create Issue `atlassian_jira_issue_screen` resources associated to different `atlassian_jira_issue_type` or projects.

```terraform
resource "atlassian_jira_issue_field_configuration" "example" {
  name = "foo"
}

resource "atlassian_jira_issue_field_configuration_item" "example_required" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id          = "customfield_10000"
    is_required = true
  }
}

resource "atlassian_jira_issue_field_configuration_item" "example_optional" {
  issue_field_configuration = atlassian_jira_issue_field_configuration.example.id
  item = {
    id          = "customfield_10001"
    is_required = false
  }
}
```

### Renderers

-> **Note** `item.renderer` can be set for text fields to the default text renderer ([`text-renderer`](https://support.atlassian.com/jira-cloud-administration/docs/configure-renderers/#Default-text-renderer)) or wiki style renderer ([`wiki-renderer`](https://support.atlassian.com/jira-cloud-administration/docs/configure-renderers/#Wiki-style-renderer)). However, `item.renderer` cannot be set for multi-select fields using the [autocomplete renderer](https://support.atlassian.com/jira-cloud-administration/docs/configure-renderers/#Autocomplete-and-select-list-renderers) or [select list renderer](https://support.atlassian.com/jira-cloud-administration/docs/configure-renderers/#Autocomplete-and-select-list-renderers). See more [details](https://support.atlassian.com/jira-cloud-administration/docs/configure-renderers/).

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `issue_field_configuration` (String) (Forces new resource) The ID of the issue field configuration.
- `item` (Attributes) Details of a field within the issue field configuration. (see [below for nested schema](#nestedatt--item))

### Read-Only

- `id` (String) The ID of the issue field configuration item. It is computed using `issue_field_configuration` and `item.id` separated by a hyphen (`-`).

<a id="nestedatt--item"></a>
### Nested Schema for `item`

Required:

- `id` (String) (Forces new resource) The ID of the field within the issue field configuration.

Optional:

- `description` (String) The description of the field within the issue field configuration.
- `is_hidden` (Boolean) Whether the field is hidden in the issue field configuration. Can be `true` or `false`.
- `is_required` (Boolean) Whether the field is required in the issue field configuration. Can be `true` or `false`.
- `renderer` (String) The renderer type for the field within the issue field configuration. Can be `text-renderer` or `wiki-renderer`.

## Import

`atlassian_jira_issue_field_configuration_item` can be imported using the `issue_field_configuration` and `item.id` separated by a comma (`,`) e.g.,

```sh
$ terraform import atlassian_jira_issue_field_configuration_item.foo 10000,customfield_10000
```
