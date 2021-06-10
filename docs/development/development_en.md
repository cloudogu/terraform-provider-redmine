# Developing terraform-provider-redmine

## Building

### Tools

Building requires these installed tools. In general, this project does not rely on cutting-edge or experimental
versions. Decent versions are okay, but if building fails you should check if you have quite outdated versions:

- Terraform client
    - f. e. v0.14.3 (linux/amd64)
- Make
    - f. e. GNU Make 4.2.1
- Docker
    - f. e. Client/server version 19.03.5
- docker-compose
    - f. e. 1.25.5
- Golang compiler
    - f. e. 1.14.12

### Local Building and other interesting `make` targets

make target | action
------------|-------
`make` / `make compile` | builds the provider binary in Go container
`make install-local` | Terraform can find the built provider when it was installed in a place known to Terraform.
`make start-redmine` | starts Redmine container into a known configuration via `docker-compose`
`make clean-redmine` | removes Redmine container
`make acceptance-test-local` | starts acceptance tests against Redmine container with local Go compiler
`make fetch-api-token` | create/refresh Redmine API token for user `admin` (is included in `acceptance-test/-local`)

### Manual Testing

Look at `examples/main.tf` and modify the `redmine_*` resources to your liking.

```bash
make install-local
make start-redmine
make fetch-api-token

cd example

terraform init # only for the first time
terraform plan
terraform apply
```

## Inner Workings

### Vision of Architecture

The usual CRUD operations for projects and issues are supported. Terraform's `import` operation is unsupported yet easy
to implement.

**Projects:**

For Redmine projects these entity fields are currently supported:

- `identifier`
- `name`
- `description`
- `homepage`
- `is_public`
- `inherit_members`

**Issues:**

For Redmine issues these entity fields are currently supported:

- `project_id` -> reference a project via Terraform resource reference
    - f. i. `redmine_project.yourtfproject.id`
- `tracker_id` -> issue trackers are currently hardcoded
    - default configuration contains these trackers:
        - Tracker ID 1 -> Bug
        - Tracker ID 2 -> Feature
        - Tracker ID 3 -> Support
- `subject` -> the title of the issue 
- `description` -> a multiline description that makes the body of the issue
- `priority_id` -> reference an issue priority from the issue priority enumeration
    - see also `local.redmine_default_issue_priorities_ids.immediate` in `examples/main.tf`
- `category_id` -> reference an issue category entity 
    - f. i. `redmine_issue_category.your_issue_category.id`

**Issue Categories:**

For Redmine issue categories these entity fields are currently supported:
- `project_id` -> reference a project via Terraform resource reference
    - f. i. `redmine_project.yourtfproject.id`
- `name` -> name of the issue category

**Versions:**

For Redmine versions these entity fields are currently supported:

- `project_id` -> reference a project via Terraform resource reference
    - f. i. `redmine_project.yourtfproject.id`
- `name` -> the human-readable name of the version
- `description` -> a single-line description of the version
- `status` -> one of these strings (defaults to "open" on creation):
  - `open`
  - `locked`
  - `closed`
- `due_date` -> the date when the version is due in the format `YYYY-MM-DD`


### Notable Aspects of Terraform Providers

Building a Terraform provider has some unexpected quirks. For instance, `*schema.ResourceData.Id()/getId()` always
points to a string ID which in turn leads to boilerplate code casting int IDs to string or leaving the ID out of an
entity if unset.

IDs in general are also a thing that are hard to test against in acceptance tests. This is because already existing data
or entity insertion sequence may lead to different IDs than expected.

### Notable changes in the used Redmine configuration

The file `docker-compose.yml` sets up Redmine without database. Redmine will default to the internal DB
provider `SQLite`. This reduces the amount of internal container wiring and avoids DB set-up scripts.

- Import Redmine default configuration
    - this includes predefined statuses and other enumerations
- Include patched settings.yml
    - only API relevant fields are changed so that Redmine is able to immediately interact via REST API
- Install SQLite in order to modify default admin user (removes barriers towards getting an API token and simplifies
  sign-on)

# Terraform Providers To-Go

Terraform cannot find local Terraform providers like this in Hashicorp's Terraform provider registry. This scenario is very likely if this provider should be used on a different machine. Anyhow, there is a convenient solution to transport this provider to another location, f. i. a Cloudogu EcoSystem.

## Packing a Terraform executable and this provider

1. checkout the Terraform Git repository and build the bundler
    - `git clone --depth 1 https://github.com/hashicorp/terraform.git ; cd terraform`
    - `go install ./tools/terraform-bundle`
    - this should leave the `terraform-bundle` binary in your `$PATH`
1. create a bundle config:

```hcl
terraform {
  # Version of Terraform to include in the bundle. An exact version number is required.
  version = "0.14.8"
}

providers {
  redmine = {
    versions = ["0.1.0"] # exact version
    source = "cloudogu.com/tf/redmine"
  }
}
```

Then it's time to buckle up, buttercup, and to bundle the binary goodness into a single ZIP archive:

```bash
$ terraform-bundle package --plugin-dir=/home/youruser/.terraform.d/plugins bundle.config
Fetching Terraform 0.14.8 core package...
Local plugin directory "/home/youruser/.terraform.d/plugins" found; scanning for provider binaries.
Found provider "cloudogu.com/tf/redmine" in "/home/youruser/.terraform.d/plugins".
- Finding cloudogu.com/tf/redmine versions matching "0.1.0"...
- Installing cloudogu.com/tf/redmine v0.1.0...
Creating terraform_0.14.8-bundle2021031713_linux_amd64.zip ...
All done!
```

## Unpacking and execution on the target machine

1. go to your favorite execution directory and unzip the Terraform package
    - `mkdir -p tf-redmine ; cd tf-redmine ; unzip terraform_0.14.8-bundle2021031713_linux_amd64.zip`
1. have a `main.tf` ready (the usual one; there is no further customizing to be done)
1. create a Terraform RC file (f. ex. `myterraformrc`)

```hcl
provider_installation {
  filesystem_mirror {
    path = "/your/dir/here/tf-redmine/providers"
  }
}
```   

Execute the provider (almost as usual) like this:

```bash
TF_CLI_CONFIG_FILE=myterraformrc ./terraform init
```

## Update documentation for Terraform

There is the official tool `tfplugindocs` from Terraform, which generates a unified documentation for the providers.
This should be updated before each release if basic things have changed at the provider.

**Note: Running the tool will remove the docs folder and recreate it. It is therefore useful to remove not generated docs before -> generate docs -> add not generated docs again.

``bash
tfplugindocs
```