FROM golang:1.23-alpine3.20 AS builder
ARG GO_BUILD_ARGS=""

WORKDIR /build
COPY ./libraries ./libraries

COPY ./go.mod ./go.sum ./
RUN go mod download
WORKDIR /build/services/order/app
COPY ./services/order/app ./
RUN go build ${GO_BUILD_ARGS} -o app

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /build/services/order/app/app ./
CMD ["./app"]