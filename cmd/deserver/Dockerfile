FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/deserver ./cmd/deserver

FROM scratch
COPY --from=builder /bin/deserver /bin/deserver
EXPOSE 8080
CMD ["/bin/deserver"]
