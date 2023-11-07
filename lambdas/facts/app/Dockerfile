FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download
WORKDIR /build/app
COPY ./ ./
RUN go build -tags lambda.norpc -o main

FROM alpine:3.18
COPY --from=builder /build/app/main ./main
CMD ["./main"]