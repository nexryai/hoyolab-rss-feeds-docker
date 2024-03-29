package main

import (
	"fmt"
	"lab.sda1.net/nexryai/hoyofeed/cache"
	"lab.sda1.net/nexryai/hoyofeed/feed"
	"lab.sda1.net/nexryai/hoyofeed/logger"
	"lab.sda1.net/nexryai/hoyofeed/server"
	"os"
	"time"
)

func main() {
	log := logger.GetLogger("Main")

	log.ProgressInfo("Init memory...")
	memory := cache.MultiTypeFeedCache{
		StarRailXml:  new(cache.MultiLangFeedCache),
		StarRailJson: new(cache.MultiLangFeedCache),
		GenshinXml:   new(cache.MultiLangFeedCache),
		GenshinJson:  new(cache.MultiLangFeedCache),
	}
	log.ProgressOk()

	// 20分おきにフィードを更新
	log.ProgressInfo("Starting feed generate goroutine...")
	go func() {
		for {
			err := feed.GenerateFeed(&memory, os.Getenv("FEED_LANG"))
			if err != nil {
				fmt.Println("err:", err)
			}

			log.Info("feed updated")

			// 20分待つ
			time.Sleep(time.Minute * 20)
		}
	}()
	log.ProgressOk()

	fmt.Println("\n")
	log.Info("starting server...")
	server.StartServer(&memory)
}
