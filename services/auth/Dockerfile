FROM golang:1.16 as builder

WORKDIR /go/src/flahmingo

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY utils utils

COPY services/auth services/auth
WORKDIR /go/src/flahmingo/services/auth
ENV CGO_ENABLED=0
RUN go build -o auth

FROM alpine:latest
USER root
WORKDIR /flahmingo

COPY --from=builder /go/src/flahmingo/services/auth/auth .
COPY setup/config.toml /etc/flahmingo/config.toml

CMD ["./auth"]
