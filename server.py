import datetime
import threading
import time

import falcon
from hoyolabrssfeeds import FeedConfigLoader, GameFeedCollection


async def generate_feeds():
    loader = FeedConfigLoader("./config.toml")

    # all games in config
    all_configs = await loader.get_all_feed_configs()
    feed_collection = GameFeedCollection.from_configs(all_configs)
    await feed_collection.create_feeds()


def generate_feeds_cron():
    while True:
        print(f"[{datetime.datetime.now()}] Generating feeds...")
        generate_feeds()
        time.sleep(20 * 60)  # Sleep for 20 minutes


class StaticResource:
    def on_get(self, req, resp, filename):
        resp.status = falcon.HTTP_200

        if filename == "genshin.json":
            resp.content_type = "application/json"
            with open("./genshin.json", 'r') as file:
                resp.text = file.read()

        elif filename == "genshin.xml":
            resp.content_type = "application/rss+xml"
            with open("./genshin.xml", 'r') as file:
                resp.text = file.read()

        elif filename == "starrail.json":
            resp.content_type = "application/json"
            with open("./starrail.json", 'r') as file:
                resp.text = file.read()

        elif filename == "starrail.xml":
            resp.content_type = "application/rss+xml"
            with open("./starrail.xml", 'r') as file:
                resp.text = file.read()

        else:
            resp.status = falcon.HTTP_NOT_FOUND
            return


app = falcon.App()
app.add_route("/ja/{filename}", StaticResource())

if __name__ == '__main__':
    from gunicorn.app.base import BaseApplication

    thread = threading.Thread(target=generate_feeds_cron)
    thread.daemon = True
    thread.start()

    class StandaloneApplication(BaseApplication):
        def __init__(self, app, options=None):
            self.options = options or {}
            self.application = app
            super().__init__()

        def load_config(self):
            config = {key: value for key, value in self.options.items() if key in self.cfg.settings and value is not None}
            for key, value in config.items():
                self.cfg.set(key.lower(), value)

        def load(self):
            return self.application

    options = {
        'bind': '0.0.0.0:8000',
        'workers': 4,
    }

    StandaloneApplication(app, options).run()
