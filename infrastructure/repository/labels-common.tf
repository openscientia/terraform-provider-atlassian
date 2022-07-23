variable "common_labels" {
  default = {
    # In alphabetical order:

    "breaking-change" = {
      color       = "b60205",
      description = "Introduces a breaking change in current functionality; usually deferred to the next major release."
    },
    "bug" = {
      color       = "b60205",
      description = "Addresses a defect in current functionality."
    },
    "dependencies" = {
      color       = "0366d6",
      description = "Used to indicate dependency changes."
    },
    "documentation" = {
      color       = "0075ca",
      description = "Introduces or discusses updates to documentation."
    },
    "enhancement" = {
      color       = "0e8a16",
      description = "Requests to existing resources that expand the functionality or scope."
    },
    "examples" = {
      color       = "63d0ff",
      description = "Introduces or discusses updates to examples."
    },
    "generators" = {
      color       = "60dea9",
      description = "Relates to code generators."
    },
    "good first issue" = {
      color       = "41e9d3",
      description = "Call to action for new contributors looking for a place to start. Smaller or straightforward issues."
    },
    "linter" = {
      color       = "0075ca",
      description = "Pertains to changes to or issues with the various linters."
    },
    "needs-triage" = {
      color       = "e99695",
      description = "Waiting for first response or review from a maintainer."
    },
    "new-data-source" = {
      color       = "7057ff",
      description = "Introduces a new data source."
    },
    "new-resource" = {
      color       = "7057ff",
      description = "Introduces a new resource."
    },
    "new-product" = {
      color       = "7057ff",
      description = "Introduces a new product."
    },
    "provider" = {
      color       = "844fba",
      description = "Pertains to the provider itself, rather than any interaction with AWS.",
    },
    "repository" = {
      color       = "587879",
      description = "Repository modifications; GitHub Actions, developer docs, issue templates, codeowners, changelog."
    },
    "technical-debt" = {
      color       = "d1ebff",
      description = "Addresses areas of the codebase that need refactoring or redesign."
    },
    "tests" = {
      color       = "60dea9",
      description = "PRs: expanded test coverage. Issues: expanded coverage, enhancements to test infrastructure."
    },

  }
  description = "Name-color-description mapping of workflow issues."
  type        = map(any)
}

resource "github_issue_label" "common" {
  for_each = var.common_labels

  repository  = "terraform-provider-atlassian"
  name        = each.key
  color       = each.value.color
  description = each.value.description
}
