FROM golang:1.20.3-alpine AS base

FROM base as dev

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install necessary packages including Git and OpenSSH
RUN apk update && apk add --no-cache ca-certificates git openssh openssh-client openssl curl && apk del libressl

RUN update-ca-certificates --fresh

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /app
CMD ["air"]

FROM base as built

WORKDIR /go/app/api

ENV CGO_ENABLED=0

RUN go get -d -v ./.. .
RUN go build -o /tmp/webscrapper_server ./*.go

FROM busybox

COPY --from=built /tmp/webscrapper_server /usr/bin/webscrapper_server
CMD ["webscrapper_server", "start"]