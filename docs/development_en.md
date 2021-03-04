# Developing terraform-provider-redmine

## Building

### Tools

Building requires these installed tools. In general, this project does not rely on cutting-edge or experimental versions. Decent versions are okay, but if building fails you should check if you have quite outdated versions:

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

The usual CRUD operations for projects and issues are supported. Terraform's `import` operation is unsupported yet easy to implement.

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
- `subject`
- `description`

### Notable Aspects of Terraform Providers

Building a Terraform provider has some unexpected quirks. For instance, `*schema.ResourceData.Id()/getId()` always points to a string ID which in turn leads to boilerplate code casting int IDs to string or leaving the ID out of an entity if unset.

IDs in general are also a thing that are hard to test against in acceptance tests. This is because already existing data or entity insertion sequence may lead to different IDs than expected.

To avoid having to manually adjust new authentication tokens with each new Redmine instance, the make target `make fetch-api-token` creates the file `examples/api_token.auto.tfvars`. This file will be read by Terraform during execution. All variables (including the variable `apikey`) must be declared to Terraform. This declaration is done in the file `examples/variables.tf`.

### Notable changes in the used Redmine configuration

The file `docker-compose.yml` sets up Redmine without database. Redmine will default to the internal DB provider `SQLite`. This reduces the amount of internal container wiring and avoids DB set-up scripts.

- Import Redmine default configuration
    - this includes predefined statuses and other enumerations
- Include patched settings.yml
    - only API relevant fields are changed so that Redmine is able to immediately interact via REST API
- Install SQLite in order to modify default admin user (removes barriers towards getting an API token and simplifies sign-on) 