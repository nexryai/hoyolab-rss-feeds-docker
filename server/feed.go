package server

import (
	"github.com/gofiber/fiber/v2"
	"lab.sda1.net/nexryai/hoyofeed/cache"
	"lab.sda1.net/nexryai/hoyofeed/logger"
	"os"
	"time"
)

func StartServer(cahce *cache.MultiTypeFeedCache) {
	app := fiber.New()
	log := logger.GetLogger("Serv")

	app.Get("/:lang/:feedType", func(ctx *fiber.Ctx) error {
		lang := ctx.Params("lang")
		feedType := ctx.Params("feedType")

		if lang == "ja" {
			lang = "ja-jp"
		}

		if lang == "" || feedType == "" {
			ctx.Status(400)
			return ctx.SendString("bad request")
		}

		multiLangCache, contentType := cache.FeedTypeToMultiLangCache(feedType, cahce)
		if multiLangCache == nil {
			ctx.Status(404)
			return ctx.SendString("not found")
		}

		feedCache := cache.LangToFeedCache(lang, multiLangCache)
		if feedCache == nil {
			ctx.Status(404)
			return ctx.SendString("not found (lang)")
		} else if feedCache.ContentBuffer == nil {
			ctx.Status(503)
			return ctx.SendString("service unavailable")
		}

		if feedCache.IsLocked {
			c := 0
			for {
				c += 1
				time.Sleep(time.Millisecond * 100)
				if !feedCache.IsLocked {
					break
				}

				if c > 100 {
					log.Warn("Too long time to wait for feedCache.IsLocked")
					ctx.Status(500)
					return ctx.SendString("internal server error")
				}
			}
		}

		ctx.Set("Content-Type", contentType)
		return ctx.Send(*feedCache.ContentBuffer)
	})

	log.Info("listening on :3000")

	err := app.Listen(":3000")
	if err != nil {
		log.FatalWithDetail("err:", err)
		os.Exit(1)
	}
}
