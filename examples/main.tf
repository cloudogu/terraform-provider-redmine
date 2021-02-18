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
  api_key = "1234567890"
}

/*
provider "redmine" {
  url = "https://192.168.56.2/redmine"
  username = "admin"
  password = "admin123"
  skip_cert_verify = true
}
*/

resource "redmine_project" "project1" {
  name = "example project"
  identifier = "example project"
  description = "this is an example project"
  is_public = false
  inherit_members = true
  tracker_ids = [ "1", "2" ]
  enabled_module_names = [ "issue_tracking", "time_tracking" ]
}

//resource "redmine_issue" "issue1" {
//  project_id = redmine_project.project1.id
//  tracker_id = 1
//  status_id = 1
//  subject = "Something should be done"
//  description = "In this ticket an **important task** should be done!\n\nGo ahead!"
//  priority_id = 1
//}
