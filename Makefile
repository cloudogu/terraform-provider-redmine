MAKEFILES_VERSION=4.3.1
OS_ARCH=linux_amd64
VERSION=1.0
ARTIFACT_ID=terraform-provider-redmine_${VERSION}_${OS_ARCH}

HOSTNAME=cloudogu.com
NAMESPACE=tf
NAME=redmine

include build/make/variables.mk

DEFAULT_ADMIN_CREDENTIALS=admin:admin
REDMINE_URL=http://localhost:8080
REDMINE_API_TOKEN=${TARGET_DIR}/redmineAPIToken.txt

TEST?=$$(go list ./... | grep -v 'vendor')

.DEFAULT_GOAL:=compile
ADDITIONAL_CLEAN=clean-test-cache

include build/make/info.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/test-integration.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/self-update.mk

install-local: compile
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: package
package: CHANGELOG.md LICENSE README.md $(BINARY)
	tar czf $(BINARY).tar.gz CHANGELOG.md LICENSE README.md $(BINARY)

PRE_INTEGRATIONTESTS=start-local-docker-compose wait-for-redmine load-redmine-defaults

.PHONY: clean-test-cache
clean-test-cache:
	@echo clean go testcache
	@go clean -testcache

.PHONY: acceptance-test
acceptance-test:
	@mkdir -p $(TARGET_DIR)/acceptance-tests
	TF_ACC=1 REDMINE_USERNAME=admin REDMINE_PASSWORD=admin \
	go test $(TEST) -coverprofile=$(TARGET_DIR)/acceptance-tests/coverage.out -timeout 120m

.PHONY: acceptance-test-local
acceptance-test-local: start-local-docker-compose
	@make acceptance-test

.PHONY: wait-for-redmine
wait-for-redmine:
	@echo "Waiting for Redmine to get up"
	@for counter in `seq 0 5`; do \
		if curl -f -s -H "Content-Type: application/json" ${REDMINE_URL}/projects.json > /dev/null; then \
			echo "Redmine is up"; \
			exit 0; \
		fi; \
		echo "waiting..."; \
		sleep 10; \
	done; \
	echo "Redmine does not seem to be up"; \
	exit 1;

.PHONY load-redmine-defaults:
load-redmine-defaults:
	@echo "Loading Redmine default configuration"
	@docker-compose exec redmine bash -c "RAILS_ENV=production REDMINE_LANG=en bundle exec rake redmine:load_default_data"

.PHONY mark-admin-password-as-changed:
mark-admin-password-as-changed: install-sqlite-client
	@echo "Mark admin password as already changed"
	@docker-compose exec redmine \
		sqlite3 /usr/src/redmine/sqlite/redmine.db \
		"update users set must_change_passwd=0 where id=1;"
	@echo "Restart Redmine to apply changed user data"
	@docker-compose restart redmine
	@make wait-for-redmine

.PHONY install-sqlite-client:
install-sqlite-client:
	@if ! docker-compose exec redmine sh -c "apk list sqlite | grep installed" ; then \
		echo "Installing sqlite..." ; \
		docker-compose exec redmine apk add sqlite ; \
	fi;

${REDMINE_API_TOKEN}: ${TARGET_DIR}
	@curl -f -s -H "Content-Type: application/json" -u ${DEFAULT_ADMIN_CREDENTIALS} ${REDMINE_URL}/my/account.json | jq -r .user.api_key > ${REDMINE_API_TOKEN}

.PHONY start-redmine:
start-redmine:
	@make start-local-docker-compose wait-for-redmine load-redmine-defaults mark-admin-password-as-changed ${REDMINE_API_TOKEN}

.PHONY teardown-redmine:
teardown-redmine:
	@docker-compose rm --force --stop