# Multi-stage build for GoDE server
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o deserver ./cmd/deserver

# Final stage
FROM alpine:3.19

# Install CA certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/deserver .

# Create non-root user
RUN addgroup -g 1000 gode && \
    adduser -D -u 1000 -G gode gode && \
    chown -R gode:gode /app

USER gode

# Expose gRPC and HTTP ports
EXPOSE 3030 8081

ENTRYPOINT ["./deserver"]
