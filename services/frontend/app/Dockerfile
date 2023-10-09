FROM golang:1.21-alpine3.18 AS builder
WORKDIR /build
COPY ./libraries ./libraries

COPY ./go.mod ./go.sum ./
RUN go mod download
WORKDIR /build/services/frontend/app
COPY ./services/frontend/app ./
RUN go build -o frontend

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /build/services/frontend/app/frontend ./
CMD ["./frontend"]