terraform {
  backend "remote" {
    organization = "openscientia"

    workspaces {
      name = "terraform-provider-atlassian-repository"
    }
  }

  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 4.0"
    }
  }

  required_version = ">= 0.13.5"
}

provider "github" {
  owner = "openscientia"
}
