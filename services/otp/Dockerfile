FROM golang:1.16 as builder

WORKDIR /go/src/flahmingo

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY utils utils

COPY services/otp services/otp
WORKDIR /go/src/flahmingo/services/otp
ENV CGO_ENABLED=0

RUN go build

FROM alpine:latest

WORKDIR /flahmingo

COPY --from=builder /go/src/flahmingo/services/otp/otp .

COPY setup/config.toml /etc/flahmingo/config.toml
CMD ["./otp"]
