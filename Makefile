BINARY_NAME=terraform-provider-sealedsecrets
BUILD_PATH=build
VERSION?=v0.1.0

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=gotestsum
GO_GET=$(GO_CMD) get

all: clean build

build: 
	$(GO_BUILD) -o $(BINARY_NAME)_v${VERSION} -v

test: 
	$(GO_GET) gotest.tools/gotestsum
	$(GO_TEST) --format short-verbose

clean: 
	$(GO_CLEAN)
	rm -rf $(BINARY_NAME)*

tidy:
	$(GO_CMD) mod tidy