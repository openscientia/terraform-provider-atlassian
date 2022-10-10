## 0.2.0 (Unreleased)

NOTES:

* data-source/atlassian_jira_issue_type_scheme: Missing configuration example has been added to documentation ([#39](https://github.com/openscientia/terraform-provider-atlassian/issues/39))
* resource/atlassian_jira_issue_type: Added [randomly generated names](https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/acctest/random.go) to execute acceptance tests in [parallel](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2@v2.21.0/helper/resource#ParallelTest) ([#98](https://github.com/openscientia/terraform-provider-atlassian/issues/98))
* resource/atlassian_jira_issue_type: Enhance acceptance tests to validate configuration changes on all attributes ([#100](https://github.com/openscientia/terraform-provider-atlassian/issues/100))
* resource/atlassian_jira_issue_type: Enhance logging output to improve provider debugging ([#101](https://github.com/openscientia/terraform-provider-atlassian/issues/101))
* resource/atlassian_jira_issue_type: Removed unnecessary [UseStateForUnknown](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework@v0.11.1/resource#UseStateForUnknown) plan modifier used with attribute [`name`](https://registry.terraform.io/providers/openscientia/atlassian/latest/docs/resources/jira_issue_type#name) ([#93](https://github.com/openscientia/terraform-provider-atlassian/issues/93))
* resource/atlassian_jira_issue_type: Replace attribute validators with [`terraform-plugin-framework-validators`](https://github.com/hashicorp/terraform-plugin-framework-validators) ([#94](https://github.com/openscientia/terraform-provider-atlassian/issues/94))
* resource/atlassian_jira_issue_type: Updated documentation to define inline `terraform import` command using markdown template. ([#53](https://github.com/openscientia/terraform-provider-atlassian/issues/53))
* resource/atlassian_jira_issue_type: Updated documentation to fix invalid link in [`Example Usage`](https://registry.terraform.io/providers/openscientia/atlassian/latest/docs/resources/jira_issue_type#example-usage) section ([#74](https://github.com/openscientia/terraform-provider-atlassian/issues/74))
* resource/atlassian_jira_issue_type_scheme: Added [randomly generated names](https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/acctest/random.go) to execute acceptance tests in [parallel](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2@v2.21.0/helper/resource#ParallelTest) ([#98](https://github.com/openscientia/terraform-provider-atlassian/issues/98))
* resource/atlassian_jira_issue_type_scheme: Enhance acceptance tests to validate configuration changes on all attributes ([#100](https://github.com/openscientia/terraform-provider-atlassian/issues/100))
* resource/atlassian_jira_issue_type_scheme: Enhance logging output to improve provider debugging ([#101](https://github.com/openscientia/terraform-provider-atlassian/issues/101))
* resource/atlassian_jira_issue_type_scheme: Replace attribute validators with [`terraform-plugin-framework-validators`](https://github.com/hashicorp/terraform-plugin-framework-validators) ([#94](https://github.com/openscientia/terraform-provider-atlassian/issues/94))
* resource/atlassian_jira_issue_type_scheme: Updated documentation to define inline `terraform import` command using markdown template. ([#53](https://github.com/openscientia/terraform-provider-atlassian/issues/53))

FEATURES:

* **New Data Source:** `atlassian_jira_issue_field_configuration` ([#63](https://github.com/openscientia/terraform-provider-atlassian/issues/63))
* **New Data Source:** `atlassian_jira_issue_field_configuration_scheme` ([#86](https://github.com/openscientia/terraform-provider-atlassian/issues/86))
* **New Data Source:** `atlassian_jira_issue_screen` ([#20](https://github.com/openscientia/terraform-provider-atlassian/issues/20))
* **New Data Source:** `atlassian_jira_issue_type_screen_scheme` ([#55](https://github.com/openscientia/terraform-provider-atlassian/issues/55))
* **New Data Source:** `atlassian_jira_myself` ([#141](https://github.com/openscientia/terraform-provider-atlassian/issues/141))
* **New Data Source:** `atlassian_jira_permission_grant` ([#122](https://github.com/openscientia/terraform-provider-atlassian/issues/122))
* **New Data Source:** `atlassian_jira_permission_scheme` ([#107](https://github.com/openscientia/terraform-provider-atlassian/issues/107))
* **New Data Source:** `atlassian_jira_project_category` ([#116](https://github.com/openscientia/terraform-provider-atlassian/issues/116))
* **New Data Source:** `atlassian_jira_screen_scheme` ([#51](https://github.com/openscientia/terraform-provider-atlassian/issues/51))
* **New Resource:** `atlassian_jira_issue_field_configuration` ([#62](https://github.com/openscientia/terraform-provider-atlassian/issues/62))
* **New Resource:** `atlassian_jira_issue_field_configuration_item` ([#73](https://github.com/openscientia/terraform-provider-atlassian/issues/73))
* **New Resource:** `atlassian_jira_issue_field_configuration_scheme` ([#78](https://github.com/openscientia/terraform-provider-atlassian/issues/78))
* **New Resource:** `atlassian_jira_issue_field_configuration_scheme_mapping` ([#90](https://github.com/openscientia/terraform-provider-atlassian/issues/90))
* **New Resource:** `atlassian_jira_issue_screen` ([#15](https://github.com/openscientia/terraform-provider-atlassian/issues/15))
* **New Resource:** `atlassian_jira_issue_type_screen_scheme` ([#54](https://github.com/openscientia/terraform-provider-atlassian/issues/54))
* **New Resource:** `atlassian_jira_permission_grant` ([#121](https://github.com/openscientia/terraform-provider-atlassian/issues/121))
* **New Resource:** `atlassian_jira_permission_scheme` ([#106](https://github.com/openscientia/terraform-provider-atlassian/issues/106))
* **New Resource:** `atlassian_jira_project_category` ([#115](https://github.com/openscientia/terraform-provider-atlassian/issues/115))
* **New Resource:** `atlassian_jira_screen_scheme` ([#47](https://github.com/openscientia/terraform-provider-atlassian/issues/47))

ENHANCEMENTS:

* provider: Add `url` attribute validation. ([#42](https://github.com/openscientia/terraform-provider-atlassian/issues/42))

## 0.1.0 (July 16, 2022)

FEATURES:

* **New Data Source:** `atlassian_jira_issue_type` ([#3](https://github.com/openscientia/terraform-provider-atlassian/issues/3))
* **New Data Source:** `atlassian_jira_issue_type_scheme` ([#5](https://github.com/openscientia/terraform-provider-atlassian/issues/5))
* **New Resource:** `atlassian_jira_issue_type` ([#3](https://github.com/openscientia/terraform-provider-atlassian/issues/3))
* **New Resource:** `atlassian_jira_issue_type_scheme` ([#5](https://github.com/openscientia/terraform-provider-atlassian/issues/5))

## 0.0.0 (July 15, 2022)

Welcome to the Terraform ATLASSIAN Provider.

The journey begins here! :rocket:
