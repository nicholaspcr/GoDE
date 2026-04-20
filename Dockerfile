# Multi-stage build for GoDE server (legacy top-level Dockerfile).
# Prefer Dockerfile.server for new work — this file is kept for backward
# compatibility with anyone still running `docker build .`.

FROM golang:1.25-alpine3.21 AS builder

WORKDIR /build

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-w -s" \
    -o /build/bin/deserver \
    ./cmd/deserver

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /build/bin/deserver /usr/local/bin/deserver

USER nonroot:nonroot

EXPOSE 3030 8081

ENTRYPOINT ["/usr/local/bin/deserver"]
