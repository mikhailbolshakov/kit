.PHONY: dep test lint mock build vendor

export GOFLAGS=-mod=vendor

vendor:
	go mod vendor

dep:
	go mod tidy

test: ## run the tests
	@echo "running tests (skipping integration)"
	go test -count=1 ./...

test-with-coverage: ## run the tests with coverage
	@echo "running tests with coverage file creation (skipping integration)"
	go test -count=1 -coverprofile .testCoverage.txt -v ./...

test-integration: ## run the integration tests
	@echo "running integration tests"
	go test -count=1 -tags integration ./...

lint:
	go vet ./...
	go fmt -mod=vendor ./...

mock: # generate mocks
	@rm -R ./mocks 2> /dev/null; \
	mockery

build: lint # build library
	go build ./...

proto: ## generates proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./grpc/*.proto ./centrifugo/proto/*.proto

artifacts: dep vendor mock build ## builds and generates all artifacts

# CI/CD gitlab commands =================================================================================================

ci-check-mocks:
	@mv ./mocks ./mocks-init
	find . -maxdepth 1 -type d \( ! -path ./*vendor ! -name . \) -exec mockery --all --dir {} \;
	mockshash=$$(find ./mocks -type f -print0 | sort -z | xargs -r0 md5sum | awk '{print $$1}' | md5sum | awk '{print $$1}'); \
	mocksinithash=$$(find ./mocks-init -type f -print0 | sort -z | xargs -r0 md5sum | awk '{print $$1}' | md5sum | awk '{print $$1}'); \
	rm -fr ./mocks-init; \
	echo $$mockshash $$mocksinithash; \
	if ! [ "$$mockshash" = "$$mocksinithash" ] ; then \
	  echo "\033[31mMocks should be updated!\033[0m" ; \
	  exit 1 ; \
	fi

ci-check: ci-check-mocks

ci-build-mr: test-with-coverage build