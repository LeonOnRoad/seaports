FROM golang:1.15-alpine AS builder
RUN apk add --no-cache
WORKDIR /tmp/build
ADD go.mod /tmp/build
ADD . /tmp/build
RUN cd /tmp/build && CGO_ENABLED=0 go build -mod=vendor ./services/ports-service

FROM scratch
WORKDIR /

COPY --from=builder /tmp/build/ports-service /ports-service
ENTRYPOINT ["/ports-service"]
