FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download
WORKDIR /build/test/app
COPY ./test/app ./
RUN go build -tags lambda.norpc -o main

FROM alpine:3.18
COPY --from=builder /build/test/app/main ./app/main
CMD ["./main"]