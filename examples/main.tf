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
  tracker_id = 1
  subject = "Something should be done"
  description = "In this ticket an **important task** should be done!\n\nGo ahead!\n\n```bash\necho -n $PATH\n```"
}