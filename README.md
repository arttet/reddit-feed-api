# Reddit Feed API

[![Go](https://img.shields.io/badge/Go-1.18-blue.svg)](https://golang.org)
[![build](https://github.com/arttet/reddit-feed-api/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/arttet/reddit-feed-api/actions/workflows/build.yml)
[![tests](https://github.com/arttet/reddit-feed-api/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/arttet/reddit-feed-api/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/arttet/reddit-feed-api/branch/main/graph/badge.svg?token=S5a5aZsotj)](https://codecov.io/gh/arttet/reddit-feed-api)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/arttet/reddit-feed-api/blob/main/LICENSE)

## Build project

### For local development
```sh
make deps  # Installation of dependencies
make build # Build the project
```
### Docker
```sh
docker-compose build
```

## Running

### For local development
```sh
make run
```

### Docker
```sh
docker-compose up -d
```

## Stopping
```sh
docker-compose down
```

## Logs
```sh
docker container logs --tail 100 reddit-feed-api_reddit-feed-api_1
```
---

## Services

### Gateway:

It reads protobuf service definitions and generates a reverse-proxy server which translates a RESTful HTTP API into gRPC

- http://localhost:8081
    - [/v1/posts](http://localhost:8081/v1/posts) - Create new posts
    - [/v1/feed/{page_id}](http://localhost:8081/v1/feed/1) - Generate a feed of posts

### gRPC:

- http://localhost:8082

### Metrics:

Metrics gRPC Server

- http://localhost:9100/metrics
    - `http_microservice_requests_total`
    - `reddit_feed_api_feed_not_found_total`

### Status:

Service condition and its information

- http://localhost:8000
    - [/live](http://localhost:8000/live)- Layed whether the server is running
    - [/ready](http://localhost:8000/ready) - Is it ready to accept requests?
    - [/version](http://localhost:8000/version) - Version and assembly information
    - [/debug/swagger/{swagger_file}](http://localhost:8000/debug/swagger/api.swagger.json) - Show a Swagger file in debug mode
    - [/debug/pprof](http://localhost:8000/debug/pprof) - Profiles when profiling enables

### Kafka

Apache Kafka is an open-source distributed event streaming platform used by thousands of companies for high-performance data pipelines, streaming analytics, data integration, and mission-critical applications.

- http://localhost:9094

### PostgreSQL

For the convenience of working with the database, you can use the [pgcli](https://github.com/dbcli/pgcli) utility. Migrations are rolled out when the service starts. migrations are located in the **./migrations** directory and are created using the [goose](https://github.com/pressly/goose) tool.

```sh
pgcli "postgresql://docker:docker@localhost:5432/reddit_feed_api"
```

---

## UI

### Swagger UI

The Swagger UI is an open source project to visually render documentation for an API defined with the OpenAPI (Swagger) Specification

- http://localhost:8080

### Grafana:

- http://localhost:3000
    - login `admin`
    - password `ADMIN`
    - [Dashboard: Reddit Feed API](http://localhost:3000/d/QXuFMwN7z/reddit-feed-api?orgId=1&refresh=5s)

### Jaeger UI

Monitor and troubleshoot transactions in complex distributed systems.

- http://localhost:16686

### Prometheus:

Prometheus is an open-source systems monitoring and alerting toolkit

- http://localhost:9090

### Kafka UI

UI for Apache Kafka is a simple tool that makes your data flows observable, helps find and troubleshoot issues faster and deliver optimal performance. Its lightweight dashboard makes it easy to track key metrics of your Kafka clusters - Brokers, Topics, Partitions, Production, and Consumption.

- http://localhost:9001

### Graylog

Graylog is a leading centralized log management solution for capturing, storing, and enabling real-time analysis of terabytes of machine data.

- http://localhost:9000
    - login `admin`
    - password `admin`
