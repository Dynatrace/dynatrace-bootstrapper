# renovate depName=github.com/golangci/golangci-lint/v2
golang_ci_cmd_version=v2.2.1
# renovate depName=github.com/daixiang0/gci
gci_version=v0.13.6
# renovate depName=golang.org/x/tools
golang_tools_version=v0.34.0
# renovate depName=github.com/vektra/mockery
mockery_version=v3.5.0
# renovate depName=github.com/igorshubovych/markdownlint-cli
markdownlint_cli_version=v0.45.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

## Install all prerequisites
prerequisites: prerequisites/setup-go-dev-dependencies prerequisites/markdownlint

## Setup go development dependencies
prerequisites/setup-go-dev-dependencies: prerequisites/go-linting prerequisites/mockery

## Install go linters
prerequisites/go-linting: prerequisites/go-deadcode
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(golang_ci_cmd_version)
	go install github.com/daixiang0/gci@$(gci_version)
	go install golang.org/x/tools/cmd/goimports@$(golang_tools_version)
	go install github.com/bombsimon/wsl/v4/cmd...@master
	go install github.com/dkorunic/betteralign/cmd/betteralign@latest

## Install go deadcode
prerequisites/go-deadcode:
	go install golang.org/x/tools/cmd/deadcode@$(golang_tools_version)

## Install go test coverage
prerequisites/go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

## Install 'markdownlint' if it is missing
prerequisites/markdownlint:
	npm install -g --force markdownlint-cli@$(markdownlint_cli_version)

## Install verktra/mockery
prerequisites/mockery:
	go install github.com/vektra/mockery/v3@$(mockery_version)

## Install 'pre-commit' if it is missing
prerequisites/setup-pre-commit:
	cp ./.github/pre-commit ./.git/hooks/pre-commit
	chmod +x ./.git/hooks/pre-commit

