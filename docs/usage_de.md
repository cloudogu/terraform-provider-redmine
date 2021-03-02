# Wie man terraform-provider-redmine verwendet

## Vorbedingungen

1. `terraform-provider-redmine` für Terraform auffindbar machen
1. eine funktionierende Redmine-Instanz, die API-Calls per API-Token entgegennimmt
1. ein Terraform-Skript

### `terraform-provider-redmine` für Terraform auffindbar machen

Da terraform-provider-redmine aktuell nicht über die Terraform-Registry verfügbar ist, muss sie manuell im Home-Verzeichnis installiert werden (beispielhaft für die Version 1.0.0 und für ein Kompilat für die Architektur linux/amd64):

```
mkdir -p ~/.terraform.d/plugins/cloudogu.com/tf/redmine/1.0.0/linux_amd64
cp  terraform-provider-redmine_1.0_linux_amd64 ~/.terraform.d/plugins/cloudogu.com/tf/redmine/1.0.0/linux_amd64
```

### eine funktionierende Redmine-Instanz

Ein bloßer Redmine-Container reicht nicht, da i. d. R. keine Konfiguration geladen wurde und zudem API-Calls deaktiviert wurden.

Aktuell ist es nur möglich, den Redmine-Provider per API-Token (`api_key`) gegenüber Redmine zu authentifizieren. 

### Terraform-Skript

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
  api_key = yourAPIToken
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
  description = "In this ticket an **important task** should be done!\n\nGo ahead!\n\n```bash\necho -n $PATH\n```"
}
```

## Terraform-Workflow

Einmalig das Terraform-Arbeitsverzeichnis initialisieren:
```
terraform init
```

Danach kann durch Hinzufügen, Verändern oder Löschen von `resource`-Blöcken im Terraform-Skript mittels dieser Befehle auf Redmine angewendet werden:

```
terraform plan # zeigt an, was Terraform während "apply" durchführen würde
terraform apply # führt die Aktion des Terraform-Skripts gegenüber Redmine durch
```