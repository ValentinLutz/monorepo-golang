FROM golang:1.21-alpine3.18 AS builder
WORKDIR /build
COPY ./libraries ./libraries

COPY ./go.mod ./go.sum ./
RUN go mod download
WORKDIR /build/services/order/app
COPY ./services/order/app ./
RUN go build -cover -o test-app

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /build/services/order/app/test-app ./
CMD ["./test-app"]