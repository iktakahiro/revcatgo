GO_CMD=GO111MODULE=on go
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_FORMAT=gofmt
GO_IMPORTS=goimports


deps:
	$(GO_GET) -v -d ./...

update:
	$(GO_GET) -v -d -u ./...

test: deps
	$(GO_GET) "github.com/stretchr/testify"
	$(GO_TEST) -v ./... -count=1

fmt:
	$(GO_GET) "golang.org/x/tools/cmd/goimports"
	find . -type f -name '*.go' | xargs $(GO_FORMAT) -s -w -l
	find . -type f -name '*.go' | xargs $(GO_IMPORTS) --local github.com/iktakahiro/revcatgo -d -e -w
