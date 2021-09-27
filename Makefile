GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.17","$(shell printf "$(GO_VERSION_SHORT)\n1.17" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.17. Found: $(GO_VERSION_SHORT))
endif

###############################################################################

SERVICE_NAME=reddit-feed-api
SERVICE_PATH=github.com/arttet/reddit-feed-api
SERVICE_MAIN=cmd/$(SERVICE_NAME)/main.go
SERVICE_EXE=./bin/$(SERVICE_NAME)$(shell go env GOEXE)

###############################################################################

# https://github.com/bufbuild/buf/releases
BUF_VERSION=v1.0.0-rc1

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GOBIN?=$(GOPATH)/bin

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

###############################################################################

.PHONY: build
build: generate .build

.PHONY: run
run:
	go run \
		-gcflags='-m' \
		-gcflags='$(SERVICE_PATH)/internal/api=-m' \
		-gcflags='$(SERVICE_PATH)/internal/server=-m' \
		$(SERVICE_MAIN) --cfg config-dev.yml

.PHONY: test
test:
	go test -v -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print ($$3)}'

.PHONY: bench
bench:
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench ./...

.PHONY: lint
lint:
	@command -v golangci-lint 2>&1 > /dev/null || (echo "Install golangci-lint" && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(shell go env GOPATH)/bin" v1.42.1)
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

# Enable pprof: internal/server/status.go:13
.PHONY: pprof
pprof: .pprof-cpu

.PHONY: grpcui
grpcui:
	grpcui -plaintext 0.0.0.0:8082

###############################################################################

.PHONY: deps
deps: .deps

.deps:
	go env -w GO111MODULE=on

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
		curl -k -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) -o "$(GOBIN)/buf" && \
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
		-o $(SERVICE_EXE) $(SERVICE_MAIN)

###############################################################################

.pprof-cpu:
	go tool pprof http://localhost:8000/debug/pprof/profile?seconds=30

.pprof-mem:
	go tool pprof -alloc_space http://localhost:8000/debug/pprof/heap

###############################################################################
