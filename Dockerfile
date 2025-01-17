FROM golang:alpine AS builder

WORKDIR /app

COPY . ./

ENV CGO_ENABLED=0

RUN apk add --no-cache ca-certificates \
 && go build -ldflags="-s -w" -buildmode=pie -trimpath -o hoyofeed main.go

FROM chimeralinux/chimera

WORKDIR /app

COPY requirements.txt /app/
RUN apk update && apk --no-cache upgrade \
    && apk --no-cache add python python-pip \
    && pip install --break-system-packages -r requirements.txt \
    && addgroup -g 816 app \
    && adduser -u 816 -G app -D -h /app app \
    && chmod +x /app/hoyofeed \
    && chown -R app:app /app

USER app
CMD [ "tini", "--", "/app/hoyofeed" ]
