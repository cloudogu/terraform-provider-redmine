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
  name = "Example Project"
  description = "This is an example project for the Super-App development team."
  homepage = "https://cloudogu.com/"
  is_public = false
  inherit_members = true
}

resource "redmine_issue" "issue1" {
  project_id = redmine_project.project1.id
  tracker_id = local.redmine_default_issue_tracker_ids.user_story
  subject = "Something should be done"
  description = <<EOT
An **important task** _should_ *be done*!

```bash
codeblock() {
  echo -n $PATH
}
```
EOT
  // the priority cannot be deleted but it can be replaced. If not provided here, Redmine takes the default or previously used priority
  priority_id = local.redmine_default_issue_priorities_ids.immediate
  // the category can be added, replaced, and be deleted
  category_id = redmine_issue_category.issue_category_dev.id
}

resource "redmine_issue_category" "issue_category_dev" {
  project_id = redmine_project.project1.id
  name = "Product Development"
}

resource "redmine_version" "issue_category_dev" {
  project_id = redmine_project.project1.id
  name = "Sprint 2021-06"
  description = "Super-App Scrum Sprint 6 (team codename: Eagle in the jar)"
  // valid values: open (default when omitted at creation), locked, closed
  status = "locked"
  // can be empty or must match date format YYYY-MM-DD
  due_date = "2021-04-01"
}