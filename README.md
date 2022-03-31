# Reddit Feed API

[![Go](https://img.shields.io/badge/Go-1.18-blue.svg)](https://golang.org)
[![build](https://github.com/arttet/reddit-feed-api/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/arttet/reddit-feed-api/actions/workflows/build.yml)
[![tests](https://github.com/arttet/reddit-feed-api/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/arttet/reddit-feed-api/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/arttet/reddit-feed-api/branch/main/graph/badge.svg?token=S5a5aZsotj)](https://codecov.io/gh/arttet/reddit-feed-api)
[![Swagger Validator](https://img.shields.io/swagger/valid/3.0?specUrl=https%3A%2F%2Fraw.githubusercontent.com%2Farttet%2Freddit-feed-api%2Fmain%2Fapi%2Fopenapi-spec%2Fapi.swagger.json)](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/arttet/reddit-feed-api/main/api/openapi-spec/api.swagger.json)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/arttet/reddit-feed-api/blob/main/LICENSE)


## Development

### Environment Setup

1. Install the latest Go from https://go.dev/
1. Install `make` and `curl` tools.
1. Install the latest Jaeger from https://www.jaegertracing.io/
1. Install the latest Prometheus from https://prometheus.io/
1. Install Kafka 2.8.0 or later from https://kafka.apache.org/
1. Install the latest Docker from https://www.docker.com/

### Build

* For local development
```sh
make all
```

* Docker
```sh
docker-compose build
```

### Usage

* For local development
```sh
go run cmd/reddit-feed-api/main.go --cfg configs/reddit-feed-api/config-dev.yml
```

* Docker
```sh
docker-compose up -d
```


## Reddit Feed API Overview

### [REST API Handlers](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/arttet/reddit-feed-api/main/api/openapi-spec/api.swagger.json)

| Method | Path                                                  | Description              |
| :----: | :---------------------------------------------------: | :----------------------: |
| POST   | [/v1/posts](http://localhost:8081/v1/posts)           | Create new posts         |
| GET    | [/v1/feed/{page_id}](http://localhost:8081/v1/feed/1) | Generate a feed of posts |

### [REST API Server (gRPC-Gateway)](https://github.com/grpc-ecosystem/grpc-gateway)

* http://localhost:8081

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
    * [/debug/swagger/{swagger_file}](http://localhost:8000/debug/swagger/api.swagger.json) - Show a Swagger file in debug mode
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

* http://localhost:8080

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
