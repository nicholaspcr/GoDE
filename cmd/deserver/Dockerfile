FROM golang:latest
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/deserver ./cmd/deserver

FROM scratch
COPY --from=0 /bin/deserver /bin/deserver
EXPOSE 8080
CMD ["/bin/deserver"]
