GO_CMD=GO111MODULE=on go
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_INSTALL=$(GO_CMD) install
GO_FORMAT=gofumpt

deps:
	$(GO_GET) -v -d ./...

update:
	$(GO_GET) -v -d -u ./...

install:
	$(GO_INSTALL) golang.org/x/tools/cmd/goimports@latest
	${GO_INSTALL} github.com/golang/mock/mockgen@latest
	${GO_INSTALL} github.com/cweill/gotests/gotests@latest
	${GO_INSTALL} github.com/fatih/gomodifytags@latest
	${GO_INSTALL} github.com/josharian/impl@latest
	${GO_INSTALL} github.com/haya14busa/goplay/cmd/goplay@latest
	${GO_INSTALL} honnef.co/go/tools/cmd/staticcheck@latest
	${GO_INSTALL} golang.org/x/tools/gopls@latest
	${GO_INSTALL} mvdan.cc/gofumpt@latest
	$(GO_GET) github.com/stretchr/testify/

test: deps
	$(GO_GET) "github.com/stretchr/testify"
	export ENV='test'; $(GO_TEST) -v ./... -count=1 -cover

fmt: install
	find . -type f -name '*.go' | xargs $(GO_FORMAT) -w -l -extra
