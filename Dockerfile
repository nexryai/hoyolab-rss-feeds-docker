FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN apk add --no-cache git ca-certificates tini g++ build-base cmake clang \
 && go build -o /app/server

FROM python:3.12-alpine

WORKDIR /app

COPY ["requirements.txt", "config.toml", "/app/"]

RUN apk add --no-cache ca-certificates tini g++ build-base cmake clang \
 && pip install --break-system-packages -r requirements.txt \
 && apk del g++ build-base cmake clang \
 && addgroup -g 816 app \
 && adduser -u 816 -G app -D -h /app app \
 && chown -R app:app /app

COPY --from=builder /app/server /app/server

USER app
CMD [ "tini", "--", "/app/server" ]