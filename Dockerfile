FROM python:3.12-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache ca-certificates tini g++ build-base cmake clang \
 && pip install --break-system-packages -r requirements.txt \
 && apk del g++ build-base cmake clang \
 && addgroup -g 816 app \
 && adduser -u 816 -G app -D -h /app app \
 && chown -R app:app /app

USER app
CMD [ "tini", "--", "python3", "server.py" ]