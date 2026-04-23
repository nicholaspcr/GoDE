# Multi-stage Dockerfile for deserver

# Stage 1: Build the Go binary
FROM golang:1.25-alpine3.21 AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 + -ldflags "-w -s" → fully static, stripped binary.
# TARGETOS/TARGETARCH are provided by BuildKit for multi-arch builds.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-w -s" \
    -o /app/bin/deserver \
    ./cmd/deserver

# Stage 2: Distroless runtime — no shell, no package manager, minimal attack surface.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin/deserver /usr/local/bin/deserver

USER nonroot:nonroot

EXPOSE 3030 8081

ENV STORE_TYPE=postgres \
    METRICS_ENABLED=true \
    METRICS_TYPE=prometheus

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/deserver", "healthcheck"]

ENTRYPOINT ["/usr/local/bin/deserver"]
CMD ["start"]
