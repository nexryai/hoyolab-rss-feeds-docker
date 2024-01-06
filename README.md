## hoyolab-rss-feeds-docker

### `compose.yml`
```
services:
  app:
    image: docker.io/nexryai/hoyorss:latest
    restart: always
    ports:
      - 127.0.0.1:8000:8000
```

### 使い方
現在日本語のみ対応しています。

```
# Genshin impact
https://example.host/ja/genshin.xml
https://example.host/ja/genshin.json

# Honkai: StarRail
https://example.host/ja/starrail.xml
https://example.host/ja/starrail.json
```
