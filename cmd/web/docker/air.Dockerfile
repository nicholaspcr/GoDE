FROM golang:1.21.3-alpine3.18

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

ENTRYPOINT ["air", "-c", "/app/cmd/web/.air.web.toml"]
