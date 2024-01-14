FROM golang:alpine AS builder

WORKDIR /app

COPY . ./

RUN apk add --no-cache git ca-certificates tini g++ build-base cmake clang \
 && go build -ldflags="-s -w" -trimpath -o hoyofeed main.go

FROM python:3.12-alpine

WORKDIR /app

COPY ["requirements.txt", "config.toml", "/app/"]
COPY --from=builder /app/hoyofeed /app/hoyofeed

RUN apk add --no-cache ca-certificates tini cmake clang build-base \
 && pip install --break-system-packages -r requirements.txt \
 && apk del g++ build-base cmake clang \
 && addgroup -g 816 app \
 && adduser -u 816 -G app -D -h /app app \
 && chmod +x /app/hoyofeed \
 && chown -R app:app /app

USER app
CMD [ "tini", "--", "/app/hoyofeed" ]
