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

Das folgende Beispielskript kann auch aus dem Verzeichnis `examples/` bezogen werden.

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
  description = "In this ticket an **important task** should be done!\n\nGo ahead!\n\n```bash\necho -n $PATH\n```"
}
```

## Terraform-Workflow

Einmalig das Terraform-Arbeitsverzeichnis initialisieren:
```
terraform init
```

Danach kann durch Hinzufügen, Verändern oder Löschen von `resource`-Blöcken im Terraform-Skript die Konfiguration mittels dieser Befehle auf Redmine angewendet werden:

```
terraform plan # zeigt an, was Terraform während "apply" durchführen würde
terraform apply # führt die Aktion des Terraform-Skripts gegenüber Redmine durch
```

# Verhalten von ausgewählten Redmine-Entitäten

## Projekte

Projekte enthalten die Felder "ID" sowie "Identifier" und können in Redmine genau null oder einmal vorkommen. Die ID ist lediglich ein technischer Bezeichner und wird beim Anlegen eines Projekts berechnet. Abgesehen davon, dass ein Projekt von anderen Entitäten referenziert wird (z. B. die Issue-Ressource im obigen Beispiel), ist die ID nicht Teil der Definition eines Projekts innerhalb eines Terraform-Skripts .

Im Gegensatz dazu ist die Projektkennung eine menschenlesbare Zeichenfolge, die nicht automatisch berechnet werden kann. Stattdessen muss der Projektbezeichner vom Benutzer gewählt werden. Da der Projektbezeichner während der Lebensdauer eines Projekts nicht geändert werden kann, wird das Ändern des Bezeichners eines bestehenden Projekts als Fehler angesehen (technisch gesehen würde Redmine diese Änderung stillschweigend ignorieren, was einen falschen Terraform-Status hinterlassen würde). Zusammenfassend lässt sich sagen, **dass es unmöglich ist, die Kennung eines bestehenden Projekts zu ändern.**

# API-Konfiguration von Redmine

Damit dieser Anbieter funktioniert, muss in Redmine mindestens der Rest-API-Zugriff aktiviert sein. Wenn dieser Provider versucht, sich mit einer Redmine-Instanz auf einem anderen Rechner zu verbinden (dazu gehören auch virtuelle Maschinen), muss in Redmine zusätzlich die JSONP-Unterstützung aktiviert sein.