MAKEFILES_VERSION=4.3.1
OS_ARCH=linux_amd64
VERSION=0.3.0
ARTIFACT_ID=terraform-provider-redmine_${VERSION}_${OS_ARCH}

HOSTNAME=cloudogu.com
NAMESPACE=tf
NAME=redmine

include build/make/variables.mk

DEFAULT_ADMIN_CREDENTIALS=admin:admin
REDMINE_URL?=http://localhost:3000
REDMINE_CONTAINERNAME?=terraform-provider-redmine_redmine_1
ACCEPTANCE_TEST_DIR=$(TARGET_DIR)/acceptance-tests
ACCEPTANCE_TEST_LOG=${ACCEPTANCE_TEST_DIR}/acceptance.test.log
ACCEPTANCE_TEST_JUNIT=${ACCEPTANCE_TEST_LOG}.xml

.DEFAULT_GOAL:=compile
ADDITIONAL_CLEAN=clean-test-cache
TF_ADDITIONAL_ENVS=
# enable to show log.Printf statements in acceptance tests (only shown if tests fail)
#TF_ACC_ADDITIONAL_ENVS=TF_LOG=TRACE

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

install-local: $(BINARY)
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: package
package: CHANGELOG.md LICENSE README.md $(BINARY)
	tar czf $(BINARY).tar.gz CHANGELOG.md LICENSE README.md $(BINARY)

PRE_INTEGRATIONTESTS=start-redmine

.PHONY: clean-test-cache
clean-test-cache:
	@echo clean go testcache
	@go clean -testcache

acceptance-test: $(BINARY) $(ACCEPTANCE_TEST_DIR)
	@TF_ACC=1 go test -v ./... -coverprofile=$(ACCEPTANCE_TEST_DIR)/coverage.out -timeout 120m 2>&1 | tee $(ACCEPTANCE_TEST_LOG)
	@cat $(ACCEPTANCE_TEST_LOG) | go-junit-report > ${ACCEPTANCE_TEST_JUNIT}
	@if grep '^FAIL' $(ACCEPTANCE_TEST_LOG); then \
		exit 1; \
	fi

.PHONY: acceptance-test-local
acceptance-test-local: $(ACCEPTANCE_TEST_DIR)
	@${TF_ADDITIONAL_ENVS} make acceptance-test

$(ACCEPTANCE_TEST_DIR):
	@echo "create acceptance-test directory"
	@mkdir -p $(ACCEPTANCE_TEST_DIR)

.PHONY: wait-for-redmine
wait-for-redmine:
	@echo "Waiting for Redmine to get up at ${REDMINE_URL}"
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
	@docker exec ${REDMINE_CONTAINERNAME} bash -c "RAILS_ENV=production REDMINE_LANG=en bundle exec rake redmine:load_default_data"

.PHONY mark-admin-password-as-changed:
mark-admin-password-as-changed: install-sqlite-client
	@echo "Mark admin password as already changed"
	@docker exec ${REDMINE_CONTAINERNAME} \
		sqlite3 /usr/src/redmine/sqlite/redmine.db \
		"update users set must_change_passwd=0 where id=1;"
	@echo "Restart Redmine to apply changed user data"
	@docker restart ${REDMINE_CONTAINERNAME}
	@make wait-for-redmine

.PHONY install-sqlite-client:
install-sqlite-client:
	@if ! docker exec ${REDMINE_CONTAINERNAME} sh -c "apk list sqlite | grep installed" ; then \
		echo "Installing sqlite..." ; \
		docker exec ${REDMINE_CONTAINERNAME} apk add sqlite ; \
	fi;

.PHONY start-redmine:
start-redmine:
	@make start-local-docker-compose wait-for-redmine load-redmine-defaults mark-admin-password-as-changed

.PHONY clean-redmine:
clean-redmine: clean
	@docker-compose rm --force --stop