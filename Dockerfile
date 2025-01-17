FROM golang:alpine AS builder

WORKDIR /app

COPY . ./

ENV CGO_ENABLED=0

RUN apk add --no-cache ca-certificates \
 && go build -ldflags="-s -w" -buildmode=pie -trimpath -o hoyofeed main.go

FROM chimeralinux/chimera

WORKDIR /app

COPY ["requirements.txt", "config.*.toml", "/app/"]
COPY --from=builder /app/hoyofeed /app/

RUN apk update && apk --no-cache upgrade \
    && apk --no-cache add python python-pip \
    && pip install --no-cache-dir --break-system-packages -r requirements.txt \
    && chmod +x /app/hoyofeed \
    && chown -R 816:816 /app

USER 816
CMD [ "/app/hoyofeed" ]
