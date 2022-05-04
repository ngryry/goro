.PHONY: test
test:
	go test -cover ./...

.PHONY: fmt
fmt:
	gofmt -l -w .
	goimports -w .

.PHONY: lint
lint:
	golangci-lint run -v