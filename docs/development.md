# Developing terraform-provider-redmine

## Inner Workings

architecture

x read

x create

-> update

delete

## docker-compose

The file `docker-compose.yml` sets up Redmine without database. Redmine will default to the internal DB provider `SQLite`. This reduces the amount of internal container wiring and avoids DB set-up scripts.

- Import Redmine default configuration
    - this includes predefined statuses and other enumerations
- Include patched settings.yml
    - only API relevant fields are changed so that Redmine is able to immediately interact via REST API
- 