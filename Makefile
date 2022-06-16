TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=bitbucket
export GO111MODULE=on

export TESTARGS=-race -coverprofile=coverage.txt -covermode=atomic

export BITBUCKET_SERVER?=http://localhost:7990
export BITBUCKET_USERNAME=admin
export BITBUCKET_PASSWORD=admin

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	#The ulimit command is required to allow the tests to open more than the default 256 files as set on MacOS. The tests will fail without this. It must be done as one
	#command otherwise the setting is lost
	ulimit -n 1024; TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -count=1

testacc-bitbucket: fmtcheck
	@bash scripts/start-docker-compose.sh
	#The ulimit command is required to allow the tests to open more than the default 256 files as set on MacOS. The tests will fail without this. It must be done as one
	#command otherwise the setting is lost
	ulimit -n 1024; TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -count=1
	@bash scripts/stop-docker-compose.sh

bitbucket-start:
	@bash scripts/start-docker-compose.sh

bitbucket-stop:
	@bash scripts/stop-docker-compose.sh

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@bash -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@bash -c "'$(CURDIR)/scripts/errcheck.sh'"

build-binaries:
	@bash -c "'$(CURDIR)/scripts/build.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

ci-build-setup:
	sudo rm -f /usr/local/bin/docker-compose
	curl -L https://github.com/docker/compose/releases/download/1.25.4/docker-compose-`uname -s`-`uname -m` > docker-compose
	chmod +x docker-compose
	sudo mv docker-compose /usr/local/bin
	bash scripts/gogetcookie.sh

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile build-binaries ci-build-setup bitbucket-start bitbucket-stop
