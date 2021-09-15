GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.17","$(shell printf "$(GO_VERSION_SHORT)\n1.17" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.17. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

###############################################################################

# https://github.com/bufbuild/buf/releases
BUF_VERSION=v0.56.0

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GOBIN?=$(GOPATH)/bin

ifeq ("NT","$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

ifeq ("Windows","$(OS_NAME)")
OS_ARCH:=$(addsuffix .exe,$(OS_ARCH))
endif

###############################################################################

SERVICE_NAME=reddit-feed-api
SERVICE_PATH=arttet/reddit-feed-api

###############################################################################

.PHONY: build
build: generate

.PHONY: run
run:
	go run cmd/reddit-feed-api/main.go

.PHONY: test
test:
	go test -v -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print ($$3)}'

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: style
style:
	find . -iname *.proto | xargs clang-format -i

.PHONY: cover
cover:
	go tool cover -html cover.out

###############################################################################

.PHONY: deps
deps: .deps

.deps:
	@ # https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

	@ # https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

	@ # https://github.com/envoyproxy/protoc-gen-validate
	go install github.com/envoyproxy/protoc-gen-validate@latest

	@ # https://github.com/grpc-ecosystem/grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

###############################################################################

.PHONY: generate
generate: .generate

.generate:
	@command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
		mkdir -p "$(GOBIN)" && \
		curl -k -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH) -o "$(GOBIN)/buf" && \
		chmod +x "$(GOBIN)/buf")
	@PATH="$(GOBIN):$(PATH)" buf generate

	mv pkg/$(SERVICE_NAME)/github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME)/* pkg/$(SERVICE_NAME)
	rm -rf pkg/$(SERVICE_NAME)/github.com/
	cd pkg/$(SERVICE_NAME) # && ls go.mod || (go mod init github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

###############################################################################
