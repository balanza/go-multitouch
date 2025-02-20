.PHONY: fmt
fmt:
	go fmt ./...

.PHOHY: lint
lint:
	golangci-lint run 

.PHOHY: lintf
lintf:
	golangci-lint run --fix

.PHONY: test
test:
	go test -v ./...
