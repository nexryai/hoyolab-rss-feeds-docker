FROM python:3.12-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache ca-certificates tini \
 && pip install --break-system-packages -r requirements.txt \
 && groupadd -g 816 app \
 && adduser -u 816 -G app -D -h /app app \
 && chown -R app:app /app

USER app
CMD [ "tini", "--", "python3", "server.py" ]