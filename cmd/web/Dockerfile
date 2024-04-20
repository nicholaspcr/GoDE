FROM golang:latest
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/web ./cmd/web

FROM scratch
COPY --from=0 /bin/web /bin/web
WORKDIR /public
COPY --from=0 /app/cmd/web/static /public/static
EXPOSE 8080
CMD ["/bin/web"]
