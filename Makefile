GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.18","$(shell printf "$(GO_VERSION_SHORT)\n1.18" | sort -V | head -1)")
$(warning NEED GO VERSION >= 1.18. Found: $(GO_VERSION_SHORT))
endif

GITHUB_PATH=github.com/arttet/reddit-feed-api

###############################################################################

.PHONY: all
all: reqs deps gen build

.PHONY: reqs
reqs: .reqs

.PHONY: deps
deps: .deps-go

.PHONY: gen
gen: .generate-go

.PHONY: build
build:  .build

.PHONY: test
test:
	go test -v -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep -v -E '100.0%|total' || echo "OK"
	go tool cover -func cover.out | grep total | awk '{print ($$3)}'

.PHONY: bench
bench:
	go test -bench ./... -benchmem -cpuprofile cpu.out -memprofile mem.out -memprofilerate 1

.PHONY: lint
lint:
	buf lint
	golangci-lint run ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	find . -iname *.go | xargs gofmt -w
	find . -iname *.proto | xargs clang-format -i

.PHONY: cover
cover:
	go tool cover -html cover.out

.PHONY: grpcui
grpcui:
	grpcui -plaintext localhost:8082

.PHONY: image
image: .image

.PHONY: debug-image
debug-image: .debug-image

.PHONY: clean
clean:
	rm -rd ./bin/ || true
	docker rm -f $(shell docker ps -a -q) || true
	docker volume rm $(shell docker volume ls -q) || true

################################################################################

# https://github.com/bufbuild/buf/releases
BUF_VERSION=v1.0.0-rc12

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

.reqs:
	@command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
		mkdir -p "$(GO_BIN)" && \
		curl -k -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) -o "$(BUF_EXE)" && \
		chmod +x "$(BUF_EXE)")

################################################################################

.PHONY: .deps-go
.deps-go:
	go mod download
	go install github.com/bold-commerce/protoc-gen-struct-transformer@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/gogo/protobuf/protoc-gen-gogofaster@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

################################################################################

.generate-go: .generate-mock .generate-reddit-feed-api

.generate-mock:
	go generate ./...

.generate-reddit-feed-api: $(eval SERVICE_NAME := reddit-feed-api) .generate-template; \
	rm -rf pkg/reddit

.generate-template:
	@ echo $(SERVICE_NAME)
	@ $(BUF_EXE) generate
	@ cp -R pkg/$(GITHUB_PATH)/pkg/* pkg/
	@ rm -rf pkg/github.com/
	@ cd pkg/$(SERVICE_NAME) && ls go.mod || (go mod init $(GITHUB_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

################################################################################

.build: .build-reddit-feed-api

.build-reddit-feed-api: \
	$(eval SERVICE_NAME := reddit-feed-api) \
	$(eval SERVICE_MAIN := cmd/$(SERVICE_NAME)/main.go) \
	$(eval SERVICE_EXE  := ./bin/$(SERVICE_NAME)) \
	.build-template

.build-template:
	CGO_ENABLED=0 go build \
		-mod=mod \
		-tags='no_mysql no_sqlite3' \
		-ldflags=" \
			-X '$(GITHUB_PATH)/internal/config.version=$(VERSION)' \
			-X '$(GITHUB_PATH)/internal/config.commitHash=$(COMMIT_HASH)' \
		" \
		-o $(SERVICE_EXE)$(shell go env GOEXE) $(SERVICE_MAIN)

################################################################################

.image: .image-reddit-feed-api

.image-reddit-feed-api: \
	$(eval SERVICE_NAME := reddit-feed-api) \
	.image-template

.image-template:
	docker build . --file deployments/docker/$(SERVICE_NAME)/Dockerfile --tag $(SERVICE_NAME):dev

################################################################################

.debug-image: .debug-image-reddit-feed-api

.debug-image-reddit-feed-api: \
	$(eval SERVICE_NAME := reddit-feed-api) \
	.debug-image-template

.debug-image-template:
	docker build . --file deployments/docker/$(SERVICE_NAME)/Dockerfile.debug --tag $(SERVICE_NAME):debug

################################################################################
