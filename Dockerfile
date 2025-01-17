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
    && chown -R 816:816 /app \
    && chown root:root /app/config.*.toml \
    && chmod 644 /app/config.*.toml \
    && chown root:root /app/hoyofeed \
    && chmod 755 /app/hoyofeed \
    && apk del python-pip

USER 816
CMD [ "/app/hoyofeed" ]
