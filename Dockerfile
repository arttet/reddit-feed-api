ARG GITHUB_PATH=github.com/arttet/reddit-feed-api


FROM golang:1.17-alpine AS builder

RUN apk add --update make git protoc protobuf protobuf-dev curl
COPY . /home/${GITHUB_PATH}
WORKDIR /home/${GITHUB_PATH}
RUN make deps && make build


FROM alpine:latest as server

LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/reddit-feed-api .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migrations/

RUN chown root:root reddit-feed-api
CMD ["./reddit-feed-api", "--migration", "up"]
