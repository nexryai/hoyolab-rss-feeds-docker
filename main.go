package main

import (
	"fmt"
	"github.com/nexryai/watchmaker"
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

	timer := watchmaker.Timer{
		BaseInterval: time.Minute * 20,
		RunOnTheHour: true,
		Delay:        time.Second * 15,
	}

	go func() {
		for {
			err := feed.GenerateFeed(&memory, os.Getenv("FEED_LANG"))
			if err != nil {
				fmt.Println("err:", err)
			}

			log.Info("feed updated")

			// 待機
			timer.WaitForNextScheduledTime()
		}
	}()
	log.ProgressOk()

	fmt.Println("\n")
	log.Info("starting server...")
	server.StartServer(&memory)
}
