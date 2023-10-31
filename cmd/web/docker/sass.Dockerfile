FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache npm \
    && npm install --global sass

WORKDIR /data

ENTRYPOINT ["npx", "sass", "--watch", "/data/static/styles/index.scss", "/data/static/styles/index.css" ]
