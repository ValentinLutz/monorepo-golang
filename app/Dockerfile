# syntax=docker/dockerfile:1

FROM golang:1.18.5-alpine3.16 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o app

# alpine:3.15.4
FROM alpine:3.16
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/app ./
CMD ["./app"]
