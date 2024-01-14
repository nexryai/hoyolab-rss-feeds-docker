package cache

func LangToFeedCache(lang string, cache *MultiLangFeedCache) *FeedCache {
	switch lang {
	case "ja-jp":
		return &cache.JaJp
	case "en-us":
		return &cache.EnUs
	case "zh-cn":
		return &cache.ZhCn
	case "zh-tw":
		return &cache.ZhTw
	case "ko-kr":
		return &cache.KoKr
	case "de-de":
		return &cache.DeDe
	case "ru-ru":
		return &cache.RuRu
	case "vi-vn":
		return &cache.ViVn
	default:
		return nil
	}
}

func FeedTypeToMultiLangCache(feedType string, cache *MultiTypeFeedCache) (*MultiLangFeedCache, string) {
	switch feedType {
	case "starrail.xml":
		return cache.StarRailXml, "application/xml+rss"
	case "starrail.json":
		return cache.StarRailJson, "application/json"
	case "genshin.xml":
		return cache.GenshinXml, "application/xml+rss"
	case "genshin.json":
		return cache.GenshinJson, "application/json"
	default:
		return nil, ""
	}
}
