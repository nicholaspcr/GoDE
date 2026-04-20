# E2E Dockerfile — used by test/e2e/docker-compose.e2e.yaml.

FROM golang:1.25-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-w -s" \
    -o /bin/deserver \
    ./cmd/deserver

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /bin/deserver /usr/local/bin/deserver

USER nonroot:nonroot

EXPOSE 3030 8081

ENTRYPOINT ["/usr/local/bin/deserver"]
CMD ["start"]
