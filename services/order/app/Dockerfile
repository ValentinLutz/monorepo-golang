# syntax=docker/dockerfile:1

#FROM golang:1.19.0-alpine3.16 AS builder
#WORKDIR /build
#COPY go.mod go.sum ./
#RUN go mod download
#COPY ./ ./
#RUN go build -o app
#
#FROM alpine:3.16
## RUN apk --no-cache add ca-certificates
#WORKDIR /app
#COPY --from=builder /build/app ./
#CMD ["./app"]

FROM golang:1.19.0-alpine3.16 AS builder
WORKDIR /build
COPY ./libraries ./libraries

WORKDIR /build/services/order/app
COPY ./services/order/app/go.mod ./services/order/app/go.sum ./
RUN go mod download
COPY ./services/order/app ./
RUN go build -o app

FROM alpine:3.16
# RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/services/order/app/app ./
CMD ["./app"]