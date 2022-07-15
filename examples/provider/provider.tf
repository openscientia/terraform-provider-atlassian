terraform {
  required_providers {
    atlassian = {
      source  = "openscientia/atlassian"
      version = "~> 0.1.0"
    }
  }
}

provider "atlassian" {
  url      = "https://foo-bar.atlassian.net"
  username = "foo@bar.com"
  apitoken = "foo&bar123"
}

resource "atlassian_issue_type" "foo" {
  name = "bar"
}
