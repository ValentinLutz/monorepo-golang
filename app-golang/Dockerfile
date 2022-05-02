# syntax=docker/dockerfile:1

FROM golang@sha256:f94174c5262af3d8446833277aa27af400fd1a880277d43ec436df06ef3bb8ab AS builder
WORKDIR /build
COPY ./ ./
RUN go build -o app

FROM alpine@sha256:a777c9c66ba177ccfea23f2a216ff6721e78a662cd17019488c417135299cd89
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/app ./
CMD ["./app"]
