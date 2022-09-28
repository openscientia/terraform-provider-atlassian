# Generated by internal/generate/repolabels/main.go; DO NOT EDIT.
variable "resource_labels" {
  default = [
    "confluence/analytics",
    "confluence/audit",
    "confluence/content",
    "confluence/contentattachments",
    "confluence/contentbody",
    "confluence/contentchildrenanddescendants",
    "confluence/contentcomments",
    "confluence/contentlabels",
    "confluence/contentmacrobody",
    "confluence/contentpermissions",
    "confluence/contentproperties",
    "confluence/contentrestrictions",
    "confluence/contentstates",
    "confluence/contentversions",
    "confluence/contentwatches",
    "confluence/dynamicmodules",
    "confluence/experimental",
    "confluence/group",
    "confluence/inlinetasks",
    "confluence/labelinfo",
    "confluence/longrunningtask",
    "confluence/relation",
    "confluence/search",
    "confluence/settings",
    "confluence/space",
    "confluence/spacepermissions",
    "confluence/spaceproperties",
    "confluence/spacesettings",
    "confluence/template",
    "confluence/themes",
    "confluence/users",
    "jira/announcementbanner",
    "jira/applicationroles",
    "jira/appmigration",
    "jira/appproperties",
    "jira/auditrecords",
    "jira/avatars",
    "jira/dashboards",
    "jira/dynamicmodules",
    "jira/filters",
    "jira/filtersharing",
    "jira/groupanduserpicker",
    "jira/groups",
    "jira/instanceinformation",
    "jira/issueattachments",
    "jira/issuecommentproperties",
    "jira/issuecomments",
    "jira/issuecustomfieldcontexts",
    "jira/issuecustomfieldoptions",
    "jira/issuecustomfields",
    "jira/issuefieldconfigurationitems",
    "jira/issuefieldconfigurations",
    "jira/issuefieldconfigurationschememappings",
    "jira/issuefieldconfigurationschemes",
    "jira/issuefields",
    "jira/issuelinks",
    "jira/issuelinktypes",
    "jira/issuenavigatorsettings",
    "jira/issuenotificationschemes",
    "jira/issuepriorities",
    "jira/issueproperties",
    "jira/issueremotelinks",
    "jira/issueresolutions",
    "jira/issues",
    "jira/issuesearch",
    "jira/issuesecuritylevel",
    "jira/issuesecurityschemes",
    "jira/issuetypeproperties",
    "jira/issuetypes",
    "jira/issuetypeschemes",
    "jira/issuetypescreenschemes",
    "jira/issuevotes",
    "jira/issuewatchers",
    "jira/issueworklogproperties",
    "jira/issueworklogs",
    "jira/jiraexpressions",
    "jira/jirasettings",
    "jira/jql",
    "jira/labels",
    "jira/myself",
    "jira/permissiongrants",
    "jira/permissions",
    "jira/permissionschemes",
    "jira/projectavatars",
    "jira/projectcategories",
    "jira/projectcomponents",
    "jira/projectemail",
    "jira/projectfeatures",
    "jira/projectkeyandnamevalidation",
    "jira/projectpermissionschemes",
    "jira/projectproperties",
    "jira/projectroleactors",
    "jira/projectroles",
    "jira/projects",
    "jira/projecttypes",
    "jira/projectversions",
    "jira/screens",
    "jira/screenschemes",
    "jira/screentabfields",
    "jira/screentabs",
    "jira/serverinfo",
    "jira/status",
    "jira/tasks",
    "jira/timetracking",
    "jira/userproperties",
    "jira/users",
    "jira/usersearch",
    "jira/webhooks",
    "jira/workflows",
    "jira/workflowschemedrafts",
    "jira/workflowschemeprojectassociations",
    "jira/workflowschemes",
    "jira/workflowstatuscategories",
    "jira/workflowstatuses",
    "jira/workflowtransitionproperties",
    "jira/workflowtransitionrules",
  ]
  description = "Set of ATLASSIAN products' resource-specific labels."
  type        = set(string)
}

resource "github_issue_label" "resource" {
  for_each = var.resource_labels

  repository  = "terraform-provider-atlassian"
  name        = each.value
  color       = "5a4edd" # color: https://registry.terraform.io/
  description = "Issues and PRs that pertain to ${each.value} resources."
}
