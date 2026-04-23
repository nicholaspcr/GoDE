# Multi-stage Dockerfile for decli

# Stage 1: Build the Go binary
FROM golang:1.25-alpine3.21 AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-w -s" \
    -o /app/bin/decli \
    ./cmd/decli

# Stage 2: Distroless runtime.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin/decli /usr/local/bin/decli

USER nonroot:nonroot

ENV DESERVER_ADDRESS=localhost:3030

ENTRYPOINT ["/usr/local/bin/decli"]
CMD ["--help"]
