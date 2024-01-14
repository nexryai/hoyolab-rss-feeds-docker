package feed

import (
	"bytes"
	"errors"
	"lab.sda1.net/nexryai/hoyofeed/cache"
	"lab.sda1.net/nexryai/hoyofeed/logger"
	"os"
	"os/exec"
)

func GenerateFeed(cachePtr *cache.MultiTypeFeedCache, lang string) error {
	log := logger.GetLogger("Feed")

	cmd := exec.Command("hoyolab-rss-feeds", "-c", "./config.toml")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Info("generating feed...")
	err := cmd.Run()

	if stderr.String() != "" {
		log.Error("stderr: ", stderr.String())
	}

	if err != nil {
		log.ErrorWithDetail("err:", err)
		return err
	}

	starRailXml, errHsrXml := os.ReadFile("./starrail.xml")
	starRailJson, errHsrJson := os.ReadFile("./starrail.json")
	genshinXml, errGenshinXml := os.ReadFile("./genshin.xml")
	genshinJson, errGenshinJson := os.ReadFile("./genshin.json")

	err = errors.Join(errHsrXml, errHsrJson, errGenshinXml, errGenshinJson)
	if err != nil {
		log.ErrorWithDetail("err:", err)
		return err
	}

	// フィードのキャッシュを書き換え
	lang = os.Getenv("FEED_LANG")

	hsrXmlCache := cache.LangToFeedCache(lang, cachePtr.StarRailXml)
	if hsrXmlCache == nil {
		log.Fatal("invalid FEED_LANG")
		os.Exit(1)
	}

	hsrJsonCache := cache.LangToFeedCache(lang, cachePtr.StarRailJson)
	if hsrXmlCache == nil {
		log.Fatal("invalid FEED_LANG")
		os.Exit(1)
	}

	genshinXmlCache := cache.LangToFeedCache(lang, cachePtr.GenshinXml)
	if hsrXmlCache == nil {
		log.Fatal("invalid FEED_LANG")
		os.Exit(1)
	}

	genshinJsonCache := cache.LangToFeedCache(lang, cachePtr.GenshinJson)
	if hsrXmlCache == nil {
		log.Fatal("invalid FEED_LANG")
		os.Exit(1)
	}

	hsrXmlCache.IsLocked = true
	hsrXmlCache.ContentBuffer = &starRailXml
	hsrXmlCache.IsLocked = false

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
