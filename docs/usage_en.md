## How to use terraform-provider-redmine

## Preconditions

1. `terraform-provider-redmine` that is discoverable for terraform
1. working Redmine instance that accepts API calls via API token
1. Terraform script

### `terraform-provider-redmine` that is discoverable for terraform

Since terraform-provider-redmine is currently not available via the Terraform registry, it must be installed manually in
the home directory (exemplary for version 1.0.0 and for a compilation for the linux/amd64 architecture):

```
mkdir -p ~/.terraform.d/plugins/cloudogu.com/tf/redmine/1.0.0/linux_amd64
cp terraform-provider-redmine_1.0_linux_amd64 ~/.terraform.d/plugins/cloudogu.com/tf/redmine/1.0.0/linux_amd64
```

### A working Redmine instance

A mere Redmine container is not sufficient, since usually no configuration has been loaded and API calls have been
disabled.

This Redmine provider authenticates against Redmine via _Basic Authentication_ with username/password pair. These values can be configured in the `redmine` provider block (see below). For this provider to work, the `Rest API` must be enabled in Redmine. Should this provider run on a machine other than the Redmine instance, Redmine must also have additional `JSONP support` enabled.

### Terraform script

The following example script can also be obtained from the `examples/` directory. 

```terraform
terraform {
  required_providers {
    redmine = {
      source = "cloudogu.com/tf/redmine"
    }
  }
}

provider "redmine" {
  url = "http://localhost:3000"
  skip_cert_verify = true
  username = "admin"
  password = "admin"
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
  tracker_id = 1
  subject = "Something should be done"
  description = "In this ticket an **important task** should be done!\n\nGo ahead!\n\n```bash\necho -n $PATH\n``"
  priority_id = 2
  category_id = redmine_issue_category.issue_category_dev.id
}

resource "redmine_issue_category" "issue_category_dev" {
  project_id = redmine_project.project1.id
  name = "Product Development"
}
```

## Terraform workflow

Initialize the Terraform working directory once:

```
terraform init
```

After that, adding, modifying or deleting `resource` blocks in the Terraform script can be applied to Redmine using
these commands:

```
terraform plan # shows what Terraform would do during "apply"
terraform apply # performs the action of the terraform script against redmine
```

# Behaviour of selected Redmine entities

## Projects

Projects contain the fields "ID" as well as "Identifier" and can exactly zero or one time in Redmine. The ID is merely a technical identifier and will be computed upon Project creation. Referencing a project from other entities aside (f. i. the issue resource in the example above), the ID is not part of defining a project within a Terraform script . 

In contrast to that, the project identifier is a human-readable string that cannot be computed automatically. Instead, the project identifier must be chosen by the user. Because the project identifier cannot be changed during a project's lifetime, changing the identifier of an existing project will be considered an error (technically Redmine silently would ignore this change which would leave a bogus Terraform state). Quintessentially, **it is impossible to change an existing project's identifier.**

## Issues

### Multiline descriptions
The issue description is a multiline text field. As such, a Redmine issue resource can not only provide single line descriptions but multiline descriptions. There are two different ways to achieve this: 

The first variant works with inlined new-lines, like so:
```terraform
resource "redmine_issue" "issue" {
  //...
  description = "An **important task** _should_ *be done*!\n\n```bash\ncodeblock() {\n  echo -n $PATH}\n```"
}
```

The second variant works with so-called (Heredocs)[https://www.terraform.io/docs/language/expressions/strings.html] which are encapsulated between double angle brackets `<<` and two identifiers. The benefit of readability in Heredocs is at hand.  

    resource "redmine_issue" "issue" {
      //...
      description = <<EOT
    An **important task** _should_ *be done*!
    
    ```bash
    codeblock() {
      echo -n $PATH
    }
    ```
    EOT
    }

### Description always marked as changed

Redmine seems to change Unix new-lines `\n` to Windows `\r\n` when it returns the description. Upon each `terraform plan`, Terraform marks the description as changed, even when there is no visible difference.

## Issue Categories

Redmine's issue categories originally contain an optional field "assigned_to" which references a user. This field is currently not supported. 

# Redmine's API configuration

In order for this provider to work, Redmine must have at least Rest API access enabled. If this provider tries to connect against a Redmine instance on a different machine (that includes Virtual Machines) then Redmine must additionally have JSONP support enabled.