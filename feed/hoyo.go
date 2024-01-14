package feed

import (
	"bytes"
	"errors"
	"fmt"
	"lab.sda1.net/nexryai/hoyofeed/cache"
	"os"
	"os/exec"
)

func exitOnNilPtr(ptr interface{}) {
	if ptr == nil {
		fmt.Println("invalid FEED_LANG")
		os.Exit(1)
	}
}

func GenerateFeed(cachePtr *cache.MultiTypeFeedCache, lang string) error {
	cmd := exec.Command("hoyolab-rss-feeds", "-c", "./config.toml")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if stderr.String() != "" {
		fmt.Println("stderr:", stderr.String())
	}

	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	starRailXml, errHsrXml := os.ReadFile("./starrail.xml")
	starRailJson, errHsrJson := os.ReadFile("./starrail.json")
	genshinXml, errGenshinXml := os.ReadFile("./genshin.xml")
	genshinJson, errGenshinJson := os.ReadFile("./genshin.json")

	err = errors.Join(errHsrXml, errHsrJson, errGenshinXml, errGenshinJson)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	// フィードのキャッシュを書き換え
	lang = os.Getenv("FEED_LANG")

	hsrXmlCache := cache.LangToFeedCache(lang, cachePtr.StarRailXml)
	exitOnNilPtr(hsrXmlCache)

	hsrJsonCache := cache.LangToFeedCache(lang, cachePtr.StarRailJson)
	exitOnNilPtr(hsrJsonCache)

	genshinXmlCache := cache.LangToFeedCache(lang, cachePtr.GenshinXml)
	exitOnNilPtr(genshinXmlCache)

	genshinJsonCache := cache.LangToFeedCache(lang, cachePtr.GenshinJson)
	exitOnNilPtr(genshinJsonCache)

	hsrXmlCache.IsLocked = true
	hsrXmlCache.ContentBuffer = &starRailXml
	hsrJsonCache.IsLocked = false

	hsrJsonCache.IsLocked = true
	hsrJsonCache.ContentBuffer = &starRailJson
	hsrJsonCache.IsLocked = false

	genshinXmlCache.IsLocked = true
	genshinXmlCache.ContentBuffer = &genshinXml
	genshinXmlCache.IsLocked = false

	genshinJsonCache.IsLocked = true
	genshinJsonCache.ContentBuffer = &genshinJson
	genshinJsonCache.IsLocked = false

	return nil
}
