package cache

type FeedCache struct {
	IsLocked      bool    // 書き換え中かどうか
	ContentBuffer *[]byte // フィードの内容
}

type MultiLangFeedCache struct {
	JaJp FeedCache
	EnUs FeedCache
	ZhCn FeedCache
	ZhTw FeedCache
	KoKr FeedCache
	DeDe FeedCache
	RuRu FeedCache
	ViVn FeedCache
}

type MultiTypeFeedCache struct {
	StarRailXml  *MultiLangFeedCache
	StarRailJson *MultiLangFeedCache
	GenshinXml   *MultiLangFeedCache
	GenshinJson  *MultiLangFeedCache
}
