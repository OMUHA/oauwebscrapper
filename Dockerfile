FROM golang:1.22-alpine AS base

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

# Configure Git to use SSH
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"


WORKDIR /app


# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

# Build the Go project
RUN go build -o /usr/bin/webscrapper_server

# Expose the port the service listens on
EXPOSE 8282
CMD ["webscrapper_server", "start"]