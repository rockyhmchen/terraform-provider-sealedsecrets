BINARY_NAME=terraform-provider-sealedsecrets
BUILD_PATH=build
VERSION?=0.1.0

GO_CMD=go
GO_TEST=gotestsum

all: clean build

build: 
	$(GO_CMD) build -o $(BINARY_NAME)_v${VERSION} -v

build_darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO_CMD) build -o $(BINARY_NAME)_darwin_amd64_v${VERSION} -v

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(BINARY_NAME)_linux_amd64_v${VERSION} -v

test: 
	$(GO_CMD) get gotest.tools/gotestsum
	$(GO_TEST) --format short-verbose

clean: 
	$(GO_CMD) clean
	rm -rf $(BINARY_NAME)*

tidy:
	$(GO_CMD) mod tidy