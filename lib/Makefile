.PHONY: lint test
GOPATH:=$(shell go env GOPATH)

test:
	-go test -v ./... -cover -gcflags=all=-l

lint:
	golangci-lint run --timeout=10m

