terraform {
  required_providers {
    redmine = {
      source = "cloudogu.com/tf/redmine"
    }
  }
}

provider "redmine" {
  url = "http://localhost:8080"
  username = "admin"
  password = "admin"
  skip_cert_verify = true
}

/*
provider "redmine" {
  url = "https://192.168.56.2/redmine"
  username = "admin"
  password = "admin123"
  skip_cert_verify = true
}
*/

resource "redmine_project" "example_project1" {
  name = "example project"
  identifier = "example project"
  description = "this is an example project"
  is_public = false
  parent_id = ""
  inherit_members = true
  tracker_ids = [ "1", "2" ]
  enabled_module_names = [ "boards", calendar, documents, files, gantt, issue_tracking, news, repository, time_tracking, wiki ]
}

resource "redmine_issue" "ticket1" {
  id = "1"
  project_identifier = "redmine_project.example_project1"
  tracker_id = "1"
  status_id = "1"
  subject = "Something should be done"
  description = "In this ticket an **important task** should be done!\n\nGo ahead!"
  priority_id = "1"
}
