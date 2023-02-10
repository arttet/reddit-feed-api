FROM golang:1.20-alpine AS builder

ARG GITHUB_PATH=github.com/arttet/reddit-feed-api
ARG VERSION COMMIT_HASH

RUN apk add --update make curl

COPY Makefile /home/${GITHUB_PATH}/Makefile
COPY go.mod /home/${GITHUB_PATH}/go.mod
COPY go.sum /home/${GITHUB_PATH}/go.sum
COPY pkg /home/${GITHUB_PATH}/pkg

WORKDIR /home/${GITHUB_PATH}
RUN make reqs
RUN make deps

COPY . /home/${GITHUB_PATH}
RUN make build


FROM alpine:latest AS server

ARG GITHUB_PATH=github.com/arttet/reddit-feed-api
ARG SERVICE_NAME=reddit-feed-api
LABEL org.opencontainers.image.source https://${GITHUB_PATH}

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/${SERVICE_NAME} .
COPY --from=builder /home/${GITHUB_PATH}/configs/${SERVICE_NAME}/config.yml .
COPY --from=builder /home/${GITHUB_PATH}/migrations/${SERVICE_NAME}/001_init_db.sql ./migrations/

RUN chown root:root ${SERVICE_NAME}
ENTRYPOINT ["./reddit-feed-api"]
