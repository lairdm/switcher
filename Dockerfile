FROM golang:1.23.2 AS builder
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o switcher main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine:latest

RUN apk update
RUN apk add --no-cache bash curl jq ddcutil

RUN adduser \
    --disabled-password \
    --uid 1000 \
    --home /app \
    --gecos '' app \
    && chown -R app /app

# This is what is currently is on rocinante
RUN addgroup -g 139 i2c

RUN addgroup app i2c

USER app

RUN mkdir -p /app \
    && chown app /app

WORKDIR /app
COPY --from=builder /workspace/switcher .
ENTRYPOINT ["/switcher"]
