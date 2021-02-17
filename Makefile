MAKEFILES_VERSION=4.3.1
OS_ARCH=linux_amd64
VERSION=1.0
ARTIFACT_ID=terraform-provider-redmine_${VERSION}_${OS_ARCH}

HOSTNAME=cloudogu.com
NAMESPACE=tf
NAME=redmine

TEST?=$$(go list ./... | grep -v 'vendor')

.DEFAULT_GOAL:=compile
ADDITIONAL_CLEAN=clean-test-cache

include build/make/variables.mk
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

PRE_INTEGRATIONTESTS=start-local-docker-compose wait-for-redmine

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
		if curl -f -s -H "Content-Type: application/json" http://localhost:8080/projects.json > /dev/null; then \
			echo "Redmine is up"; \
			exit 0; \
		fi; \
		echo "waiting..."; \
		sleep 10; \
	done; \
	echo "Redmine does not seem to be up"; \
	exit 1;
