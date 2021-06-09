# terraform-provider-redmine entwickeln

## Erstellen

### Tools

Das Bauen erfordert diese installierten Werkzeuge. Im Allgemeinen verlässt sich dieses Projekt nicht auf topaktuelle oder experimentelle Versionen. Halbwegs aktuelle Versionen sind in Ordnung, aber wenn das Bauen fehlschlägt, sollte überprüft werden, ob stark veraltete Versionen der jeweiligen Tools vorliegen:

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

### Lokal Bauen und andere interessante `make`-Targets

make-Target | Aktion
------------|-------
`make` / `make compile` | baut das Provider-Binary in einem Go-Container
`make install-local` | Terraform kann den gebauten Provider finden, wenn er an einem Terraform bekannten Ort installiert wurde.
`make start-redmine` | startet den Redmine-Container in einer bekannten Konfiguration per `docker-compose`
`make clean-redmine` | entfernt den Redmine-Container
`make acceptance-test-local` | startet Akzeptanztests gegen Redmine-Container mit lokalem Go-Compiler
`make fetch-api-token` | erzeugt/aktualisiert Redmine-API-Token für Benutzer admin (ist in Target `acceptance-test/-local` enthalten) 

### Manuelles Testen

Die Beispieldatei `examples/main.tf` kann direkt genutzt werden. Die `redmine_*` Ressourcen können individuell angepasst werden. 

```bash
make install-local
make start-redmine
make fetch-api-token

cd example

terraform init # wird i. d. R. nur einmalig nötig
terraform plan
terraform apply
```

## Innenleben

### Architekturvision 

Die üblichen CRUD-Operationen für Projekte und Issues werden unterstützt. Die `import`-Operation von Terraform wird nicht unterstützt, ist aber leicht zu implementieren.

**Projects:**

Für Redmine-Projekte werden derzeit diese Objektfelder unterstützt:

- `identifier`
- `name`
- `description`
- `homepage`
- `is_public`
- `inherit_members`

**Issues:**

Für Redmine-Issues werden derzeit diese Objektfelder unterstützt:

- `project_id` -> referenziert ein Projekt über eine Terraform-Ressourcenreferenz
      - z. B. `redmine_project.yourtfproject.id`
- `tracker_id` -> Tracker sind derzeit hardcodiert 
   - Die Standardkonfiguration enthält diese Tracker:
      - Tracker ID 1 -> Bug
      - Tracker ID 2 -> Feature
      - Tracker ID 3 -> Support
- `subject` -> the title of the issue
- `description` -> a multline description that makes the body of the issue
- `priority_id` -> referenziert eine Issue-Priorität aus der Issue-Prioritäten-Aufzählung
    - siehe auch `local.redmine_default_issue_priorities_ids.immediate` in `examples/main.tf`
- `category_id` -> referenziert eine Terraform-Issue-Kategorie-Ressourcenreferenz
    - z. B. `redmine_issue_category.your_issue_category.id`

**Issue Categories:**

Für Redmine-Ticketkategorien werden derzeit diese Objektfelder unterstützt:
- `project_id` -> referenziert ein Projekt über eine Terraform-Ressourcenreferenz
  - z. B. `redmine_project.yourtfproject.id`
- `name` -> name of the issue category

**Versions:**

Für Redmine-Versionen werden derzeit diese Objektfelder unterstützt:

- `project_id` -> referenziert ein Projekt über eine Terraform-Ressourcenreferenz
  - z. B. `redmine_project.yourtfproject.id`
- `name` -> ein menschenlesbarer Name der Version
- `description` -> eine einzeilige Beschreibung der Version
- `status` -> eine dieser Zeichenketten (standardmäßig "open" bei der Erstellung):
    - `open`
    - `locked`
    - `closed`
- `due_date` -> das Datum, an dem die Version fällig ist, im Format `YYYY-MM-DD`

### Erwähnenswerte Aspekte von Terraform-Providern

Die Erstellung eines Terraform-Providers hat einige Eigenheiten. Zum Beispiel zeigt `*schema.ResourceData.Id()/getId()` **immer** auf eine String-ID, was wiederum zu Boilerplate-Code führt, in dem int-IDs in einen String umwandelt oder die ID aus einer Entität herauslässt, wenn sie nicht gesetzt ist.

IDs im Allgemeinen sind auch eine Sache, die in Akzeptanztests schwer testbar sind. Das liegt daran, dass bereits vorhandene Daten oder die Einfügereihenfolge von Fachobjekten zu anderen IDs führen können als erwartet.

### Erwähnenswerte Änderungen in der verwendeten Redmine-Konfiguration

Die Datei `docker-compose.yml` richtet Redmine ohne Datenbank ein. Redmine wird standardmäßig den internen DB-Provider `SQLite` verwenden. Dies reduziert den Umfang der internen Container-Verkabelung und vermeidet DB-Setup-Skripte.

- Redmine-Standardkonfiguration importieren
    - dies beinhaltet vordefinierte Status und andere Aufzählungen
- Gepatchte `settings.yml` einbinden
    - nur API-relevante Felder werden geändert, damit Redmine sofort über die REST-API interagieren kann
- SQLite installieren, um den Standard-Admin-Benutzer zu ändern (beseitigt die Hürden für den Erhalt eines API-Tokens und vereinfacht die Anmeldung)

# Terraform Providers To-Go

Terraform kann lokale Terraform Provider nicht in der Terraform Provider Registry von Hashicorp finden. Dieses Szenario tritt bspw. wahrscheinlich ein, wenn dieser Provider auf einem anderen Rechner verwendet werden soll. Es gibt jedoch eine bequeme Lösung, um diesen Provider an einen anderen Ort zu transportieren, z. B. in ein Cloudogu EcoSystem.

## Packen einer ausführbaren Terraform-Datei und dieses Providers

1. das Terraform Git Repository auschecken und den Bundler bauen
    - `git clone --depth 1 https://github.com/hashicorp/terraform.git ; cd terraform`
    - `go install ./tools/terraform-bundle`
    - dies sollte die `terraform-bundle`-Ausführungsdatei im `$PATH` hinterlassen
1. eine Bundle-Konfiguration anlegen:

```hcl
terraform {
  # Version von Terraform, die in das Bundle aufgenommen werden soll. Eine genaue Versionsnummer ist erforderlich.
  version = "0.14.8"
}

provider {
  redmine = {
    versions = ["0.1.0"] # genaue Version
    source = "cloudogu.com/tf/redmine"
  }
}
```

Dann ist es an der Zeit, sich anzuschnallen, die mühsam erzeugten Bits in ein eigenes ZIP-Archiv zu verpacken:

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

## Entpacken und Ausführen auf dem Zielrechner

1. in das bevorzugte Ausführungsverzeichnis gehen und das Terraform-Paket entpacken
    - `mkdir -p tf-redmine ; cd tf-redmine ; unzip terraform_0.14.8-bundle2021031713_linux_amd64.zip`
1. eine `main.tf` bereithalten (das Übliche; hier müssen keine weiteren Anpassungen vorgenommen werden)
1. eine Terraform RC-Datei anlegen (z. B. `myterraformrc`)

```hcl
provider_installation {
  filesystem_mirror {
    path = "/Ihr/Verzeichnis/hier/tf-redmine/providers"
  }
}
```   

Der Provider kann dann mit Übergabe der RC-Datei (fast wie üblich) ausgeführt werden:

```bash
TF_CLI_CONFIG_FILE=myterraformrc ./terraform init
```

## Dokumentation für Terraform aktualisieren

Es gibt das offizielle Tool `tfplugindocs` von Terraform, welches eine einheitliche Dokumentation für die Provider erzeugt.
Diese sollte vor jeden Release aktualisiert werden, wenn sich grundlegende Dinge am Provider geändert haben.

**Hinweis:** Das Ausführen von dem Tool entfernt den docs Ordner und erstellt ihn neu. Es ist daher sinnvoll nicht generierte Doku vorher herauszunehmen -> Doku zu generieren -> nicht generierte Doku wieder hinzuzufügen.

```bash
tfplugindocs
```