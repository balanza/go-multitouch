GOLANGCI_LINT_VERSION := v1.61.0
GOLANGCI_LINT ?= go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: fmt
fmt:
	go fmt ./...

.PHOHY: lint
lint:
	$(GOLANGCI_LINT) run 

.PHOHY: lintf
lintf:
	$(GOLANGCI_LINT) run --fix


.PHONY: test
test:
	go test -v ./...
