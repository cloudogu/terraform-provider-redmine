# Developing terraform-provider-redmine

## docker-compose

The file `docker-compose.yml` sets up Redmine without database. Redmine will default to the internal DB provider `SQLite`. This reduces the amount of internal container wiring and avoids DB set-up scripts.