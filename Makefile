GO ?= go

GOLANGCI ?= golangci-lint
GOVULNCHECK ?= $(GO) tool govulncheck

.PHONY: deps update tidy tidy-check fmt fmt-check lint test vulncheck

deps:
	$(GO) mod download

update:
	$(GO) get -u ./...
	$(GO) mod tidy

tidy:
	$(GO) mod tidy

tidy-check:
	$(GO) mod tidy -diff

fmt:
	$(GOLANGCI) fmt ./...

fmt-check:
	$(GOLANGCI) fmt --diff ./...

lint:
	$(GOLANGCI) run --modules-download-mode=readonly ./...

test:
	ENV='test' $(GO) test -race -cover ./... -count=1

vulncheck:
	$(GOVULNCHECK) ./...
