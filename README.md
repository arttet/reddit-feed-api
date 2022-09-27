# Reddit Feed API

[![Go](https://img.shields.io/badge/Go-1.19-blue.svg)](https://golang.org)
[![build](https://github.com/arttet/reddit-feed-api/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/arttet/reddit-feed-api/actions/workflows/build.yml)
[![tests](https://github.com/arttet/reddit-feed-api/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/arttet/reddit-feed-api/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/arttet/reddit-feed-api/branch/main/graph/badge.svg?token=S5a5aZsotj)](https://codecov.io/gh/arttet/reddit-feed-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/arttet/reddit-feed-api)](https://goreportcard.com/report/github.com/arttet/reddit-feed-api)
[![Swagger Validator](https://img.shields.io/swagger/valid/3.0?specUrl=https%3A%2F%2Fraw.githubusercontent.com%2Farttet%2Freddit-feed-api%2Fmain%2Fapi%2Fopenapi-spec%2Fapi.swagger.json)](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/arttet/reddit-feed-api/main/api/openapi-spec/api.swagger.json)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/arttet/reddit-feed-api/blob/main/LICENSE)

## Environment Setup

1. Install the latest Go from https://go.dev/
1. Install `make` and `curl` tools.
1. Install PostgreSQL 13 or later from https://www.postgresql.org/download/
1. Install the latest Jaeger from https://www.jaegertracing.io/
1. Install the latest Prometheus from https://prometheus.io/
1. Install OpenJDK 17 from https://openjdk.java.net/projects/jdk/17/
1. Install Kafka 3.1.0 or later from https://kafka.apache.org/
1. Install the latest Docker from https://www.docker.com/

## Usage

```sh
$ make help
  Local commands
help:    Show this help
all:     Run the all build commands
reqs:    Install requirements
deps:    Build and install Go dependencies
gen:     Generate Go code
build:   Compile the current package and all of its dependencies
test:    Execute the unit tests
bench:   Execute the benchmark tests
lint:    Check the current package to catch common mistakes and improve the code
tidy:    Cleanup go.mod
fmt:     Format *.go and *.proto files using gofmt and clang-format
cover:   Show the cover report
grpcui:  Run the grpcui tool
clean:   Remove generated artifacts
update:  Update dependencies as recorded in the go.mod and go.sum files
  Docker commands
image:           Build Dockerfile
debug-image:     Build Dockerfile.debug
dc-build:        Build docker-compose.yaml
up:              Up reddit-feed-api
down:            Down reddit-feed-api
rm:              Remove Docker artifacts
  CLI commands
run:             Run reddit-feed-api locally
cli-create:      Run the CreatePostsV1 handle
cli-generate:    Run the GenerateFeedV1 handle
producer:        Run the Kafka producer command
consumer:        Run the Kafka consumer command
```

## Reddit Feed API Overview

### [REST API Handlers](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/arttet/reddit-feed-api/main/api/openapi-spec/api.swagger.json)

| Method | Path                                                  | Description              |
| :----: | :---------------------------------------------------: | :----------------------: |
| POST   | [/v1/posts](http://localhost:8081/v1/posts)           | Create new posts         |
| GET    | [/v1/feed/{page_id}](http://localhost:8081/v1/feed/1) | Generate a feed of posts |

### [REST API Server (gRPC-Gateway)](https://github.com/grpc-ecosystem/grpc-gateway)

* http://localhost:8080

### [gRPC Server](https://github.com/grpc/grpc-go)

* http://localhost:8082

### Metrics Server

* http://localhost:9100/metrics
    * `http_microservice_requests_total`
    * `reddit_feed_api_feed_not_found_total`

### Status Server

* http://localhost:8000
    * [/live](http://localhost:8000/live)- Layed whether the server is running
    * [/ready](http://localhost:8000/ready) - Is it ready to accept requests?
    * [/version](http://localhost:8000/version) - Version and assembly information
    * [/debug/pprof](http://localhost:8000/debug/pprof) - Profiles when profiling enables

### [Apache Kafka](https://kafka.apache.org/)

* http://localhost:9094

### [PostgreSQL](https://www.postgresql.org/docs/)

* [pgcli](https://github.com/dbcli/pgcli)

```sh
pgcli "postgresql://docker:docker@localhost:5432/reddit_feed_api"
```

* [goose](https://github.com/pressly/goose)

```sh
goose -dir migrations/reddit-feed-api postgres "postgresql://docker:docker@localhost:5432/reddit_feed_api" status
goose -dir migrations/reddit-feed-api postgres "postgresql://docker:docker@localhost:5432/reddit_feed_api" up
goose -dir migrations/reddit-feed-api postgres "postgresql://docker:docker@localhost:5432/reddit_feed_api" status
```


## Tools (Use docker-compose)

### [Jaeger UI](https://www.jaegertracing.io/)

* http://localhost:16686

### [Prometheus UI](https://prometheus.io/)

* http://localhost:9090


### [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/arttet/reddit-feed-api/main/api/openapi-spec/api.swagger.json)

* http://localhost:7000

### [Grafana](https://grafana.com/oss/grafana/)

* http://localhost:3000
    * login `admin`
    * password `ADMIN`
    * [Dashboard: Reddit Feed API](http://localhost:3000/d/QXuFMwN7z/reddit-feed-api?orgId=1&refresh=5s)

### [Kafka UI](https://github.com/provectus/kafka-ui)

* http://localhost:9001

### [Graylog](https://www.graylog.org/)

* http://localhost:9000
    * login `admin`
    * password `admin`
