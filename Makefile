GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.17","$(shell printf "$(GO_VERSION_SHORT)\n1.17" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.17. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

###############################################################################

SERVICE_NAME=reddit-feed-api
SERVICE_PATH=github.com/arttet/reddit-feed-api

###############################################################################

# https://github.com/bufbuild/buf/releases
BUF_VERSION=v1.0.0-rc1

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GOBIN?=$(GOPATH)/bin

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

ifeq ("Windows", "$(OS_NAME)")
OS_ARCH:=$(addsuffix .exe,$(OS_ARCH))
endif

###############################################################################

.PHONY: build
build: generate .build
ifeq ("NT", "$(findstring NT,$(shell uname -s))")
	mv ./bin/$(SERVICE_NAME) ./bin/$(SERVICE_NAME).exe
endif

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

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: style
style:
	find . -iname *.go | xargs gofmt -w
	find . -iname *.proto | xargs clang-format -i

.PHONY: cover
cover:
	go tool cover -html cover.out

.PHONY: grpcui
grpcui:
	grpcui -plaintext 0.0.0.0:8082

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

	mv pkg/$(SERVICE_NAME)/$(SERVICE_PATH)/pkg/$(SERVICE_NAME)/* pkg/$(SERVICE_NAME)
	rm -rf pkg/$(SERVICE_NAME)/github.com/
	cd pkg/$(SERVICE_NAME) && ls go.mod || (go mod init $(SERVICE_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

###############################################################################

.build:
	go mod download && CGO_ENABLED=0 go build \
		-mod=mod \
		-tags='no_mysql no_sqlite3' \
		-ldflags=" \
			-X '$(SERVICE_PATH)/internal/config.version=$(VERSION)' \
			-X '$(SERVICE_PATH)/internal/config.commitHash=$(COMMIT_HASH)' \
		" \
		-o ./bin/reddit-feed-api ./cmd/reddit-feed-api/main.go
