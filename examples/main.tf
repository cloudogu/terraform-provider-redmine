terraform {
  required_providers {
    redmine = {
      source = "cloudogu.com/tf/redmine"
    }
  }
}

provider "redmine" {
  url = "http://localhost:3000"
  username = var.username
  password = var.password
  skip_cert_verify = true
  api_key = var.apikey
}

/*
provider "redmine" {
  url = "https://192.168.56.2/redmine"
  username = "admin"
  password = "admin123"
  skip_cert_verify = true
}
*/

locals {
  redmine_default_issue_tracker_ids = {
    bug = 1
    user_story = 2
    support = 3
  }

  redmine_default_issue_priorities_ids = {
    low = 1
    normal = 2
    high = 3
    urgent = 4
    immediate = 5
  }

  redmine_default_issue_status_ids = {
    new = 1
    in_progress = 2
    resolved = 3
    feedback = 4
    closed = 5
    rejected = 6
  }
}

resource "redmine_project" "project1" {
  identifier = "exampleproject"
  name = "example project"
  description = "this is an example project."
  homepage = "https://cloudogu.com/"
  is_public = false
  inherit_members = true
}

resource "redmine_issue" "issue1" {
  project_id = redmine_project.project1.id
  tracker_id = local.redmine_default_issue_tracker_ids.user_story
  subject = "Something should be done"
  description = "In this ticket an **important task** should be done!\n\nGo ahead!\n\n```bash\necho -n $PATH\n```"
  priority_id = local.redmine_default_issue_priorities_ids.immediate
}