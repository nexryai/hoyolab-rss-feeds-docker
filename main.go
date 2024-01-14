package main

import (
	"fmt"
	"lab.sda1.net/nexryai/hoyofeed/cache"
	"lab.sda1.net/nexryai/hoyofeed/feed"
	"lab.sda1.net/nexryai/hoyofeed/server"
	"os"
	"time"
)

func main() {
	memory := cache.MultiTypeFeedCache{
		StarRailXml:  new(cache.MultiLangFeedCache),
		StarRailJson: new(cache.MultiLangFeedCache),
		GenshinXml:   new(cache.MultiLangFeedCache),
		GenshinJson:  new(cache.MultiLangFeedCache),
	}

	// 20分おきにフィードを更新
	go func() {
		for {
			err := feed.GenerateFeed(&memory, os.Getenv("FEED_LANG"))
			if err != nil {
				fmt.Println("err:", err)
			}

			fmt.Println("feed updated")

			// 20分待つ
			time.Sleep(time.Minute * 10)
		}
	}()

	server.StartServer(&memory)
}
