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

Schauen Sie sich `examples/main.tf` und ändern Sie die `redmine_*` Ressourcen nach Ihren Wünschen. 

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

Für Redmine-Projekte werden derzeit diese Objektfelder unterstützt::

- `project_id` -> referenziert ein Projekt über eine Terraform-Ressourcenreferenz
      - z. B. `redmine_project.yourtfproject.id`
- `tracker_id` -> Tracker sind derzeit hardcodiert 
   - Die Standardkonfiguration enthält diese Tracker:
      - Tracker ID 1 -> Bug
      - Tracker ID 2 -> Feature
      - Tracker ID 3 -> Support
- `subject`
- `description`

### Erwähnenswerte Aspekte von Terraform-Providern

Die Erstellung eines Terraform-Providers hat einige Eigenheiten. Zum Beispiel zeigt `*schema.ResourceData.Id()/getId()` **immer** auf eine String-ID, was wiederum zu Boilerplate-Code führt, in dem int-IDs in einen String umwandelt oder die ID aus einer Entität herauslässt, wenn sie nicht gesetzt ist.

IDs im Allgemeinen sind auch eine Sache, die in Akzeptanztests schwer testbar sind. Das liegt daran, dass bereits vorhandene Daten oder die Einfügereihenfolge von Fachobjekten zu anderen IDs führen können als erwartet.

Damit nicht bei jeder neuen Redmine-Instanz neue Authentifizierungstoken händisch angepasst werden müssen, erzeugt das Make-Target `make fetch-api-token` die Datei `examples/api_token.auto.tfvars`. Diese wird von Terraform während der Ausführung eingelesen. Sämtliche Variablen (also auch die eben besprochene Variable `apikey`) müssen gegenüber Terraform deklariert werden. Diese Deklaration geschieht in der Datei `examples/variables.tf`.

### Erwähnenswerte Änderungen in der verwendeten Redmine-Konfiguration

Die Datei `docker-compose.yml` richtet Redmine ohne Datenbank ein. Redmine wird standardmäßig den internen DB-Provider `SQLite` verwenden. Dies reduziert den Umfang der internen Container-Verkabelung und vermeidet DB-Setup-Skripte.

- Redmine-Standardkonfiguration importieren
    - dies beinhaltet vordefinierte Status und andere Aufzählungen
- Gepatchte `settings.yml` einbinden
    - nur API-relevante Felder werden geändert, damit Redmine sofort über die REST-API interagieren kann
- SQLite installieren, um den Standard-Admin-Benutzer zu ändern (beseitigt die Hürden für den Erhalt eines API-Tokens und vereinfacht die Anmeldung)
