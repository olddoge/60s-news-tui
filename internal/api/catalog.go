package api

import (
	"net/url"
	"strings"
)

// EndpointMenuText holds localized display text for a supported endpoint menu item.
type EndpointMenuText struct {
	Emoji string
	Zh    string
	En    string
}

// EndpointMenuTexts is the shared endpoint menu catalog.
var EndpointMenuTexts = map[string]EndpointMenuText{
	"/v2/60s":                        {Emoji: "📰", Zh: "每天 60 秒读懂世界", En: "Understand the World in 60 Seconds"},
	"/v2/60s/rss":                    {Emoji: "📡", Zh: "rss 订阅", En: "RSS Feed"},
	"/v2/answer":                     {Emoji: "📖", Zh: "随机答案之书", En: "Random Book of Answers"},
	"/v2/baike":                      {Emoji: "📚", Zh: "百度百科词条", En: "Baidu Baike Entry"},
	"/v2/bili":                       {Emoji: "📺", Zh: "哔哩哔哩热搜", En: "Bilibili Trending Searches"},
	"/v2/bing":                       {Emoji: "🖼️", Zh: "必应每日壁纸", En: "Bing Daily Wallpaper"},
	"/v2/changya":                    {Emoji: "🎤", Zh: "随机唱歌音频", En: "Random Singing Audio"},
	"/v2/chemical":                   {Emoji: "⚗️", Zh: "", En: ""},
	"/v2/douyin":                     {Emoji: "🎵", Zh: "抖音热搜", En: "Douyin Trending Searches"},
	"/v2/duanzi":                     {Emoji: "😄", Zh: "随机搞笑段子", En: "Random Jokes"},
	"/v2/epic":                       {Emoji: "🎮", Zh: "Epic Games 游戏", En: "Epic Games"},
	"/v2/exchange-rate":              {Emoji: "💱", Zh: "当日货币汇率", En: "Daily Exchange Rates"},
	"/v2/fabing":                     {Emoji: "📝", Zh: "随机发病文学", En: "Random Melodramatic Copy"},
	"/v2/hitokoto":                   {Emoji: "💬", Zh: "随机一言", En: "Random Hitokoto Quote"},
	"/v2/ip":                         {Emoji: "🌐", Zh: "公网 IP 地址", En: "Public IP Address"},
	"/v2/kfc":                        {Emoji: "🍗", Zh: "随机 KFC 文案", En: "Random KFC Copy"},
	"/v2/luck":                       {Emoji: "🍀", Zh: "随机运势", En: "Random Fortune"},
	"/v2/today-in-history":           {Emoji: "📅", Zh: "历史上的今天", En: "Today in History"},
	"/v2/toutiao":                    {Emoji: "🔥", Zh: "头条热搜榜", En: "Toutiao Trending List"},
	"/v2/weibo":                      {Emoji: "📣", Zh: "微博热搜", En: "Weibo Trending Searches"},
	"/v2/zhihu":                      {Emoji: "💡", Zh: "知乎话题榜", En: "Zhihu Topic Rankings"},
	"/v2/lunar":                      {Emoji: "🌙", Zh: "农历信息", En: "Lunar Calendar Info"},
	"/v2/ai-news":                    {Emoji: "🤖", Zh: "AI 资讯快报", En: "AI News Briefing"},
	"/v2/it-news":                    {Emoji: "💻", Zh: "实时 IT 资讯", En: "Real-Time IT News"},
	"/v2/it-news/rank":               {Emoji: "🏆", Zh: "IT 之家热门榜单", En: "ITHome Hot Rankings"},
	"/v2/awesome-js":                 {Emoji: "🧩", Zh: "随机 JS 趣味题", En: "Random JavaScript Quiz"},
	"/v2/qrcode":                     {Emoji: "🔳", Zh: "生成二维码", En: "Generate QR Code"},
	"/v2/dad-joke":                   {Emoji: "🥶", Zh: "随机冷笑话", En: "Random Dad Joke"},
	"/v2/rednote":                    {Emoji: "📕", Zh: "小红书热点", En: "Rednote Hot Topics"},
	"/v2/dongchedi":                  {Emoji: "🚗", Zh: "懂车帝热搜", En: "Dongchedi Trending Searches"},
	"/v2/moyu":                       {Emoji: "🐟", Zh: "摸鱼日报", En: "Workday Slack-Off Daily"},
	"/v2/quark":                      {Emoji: "🔎", Zh: "夸克热点", En: "Quark Hot Topics"},
	"/v2/whois":                      {Emoji: "🪪", Zh: "Whois 查询", En: "Whois Lookup"},
	"/v2/health":                     {Emoji: "💪", Zh: "身体健康分析", En: "Health Analysis"},
	"/v2/password":                   {Emoji: "🔐", Zh: "密码生成器", En: "Password Generator"},
	"/v2/password/check":             {Emoji: "🛡️", Zh: "密码强度检测", En: "Password Strength Check"},
	"/v2/maoyan/all/movie":           {Emoji: "🎬", Zh: "猫眼全球票房总榜", En: "Maoyan Global Box Office Rankings"},
	"/v2/maoyan/realtime/movie":      {Emoji: "🎟️", Zh: "猫眼电影实时票房", En: "Maoyan Real-Time Movie Box Office"},
	"/v2/maoyan/realtime/tv":         {Emoji: "📺", Zh: "猫眼电视收视排行", En: "Maoyan Real-Time TV Ratings"},
	"/v2/maoyan/realtime/web":        {Emoji: "🌡️", Zh: "猫眼网剧实时热度", En: "Maoyan Real-Time Web Series Heat"},
	"/v2/hacker-news/new":            {Emoji: "🆕", Zh: "Hacker News 新帖", En: "Hacker News New"},
	"/v2/hacker-news/top":            {Emoji: "⬆️", Zh: "Hacker News 排行", En: "Hacker News Top"},
	"/v2/hacker-news/best":           {Emoji: "⭐", Zh: "Hacker News 最佳", En: "Hacker News Best"},
	"/v2/baidu/hot":                  {Emoji: "🔥", Zh: "百度实时热搜", En: "Baidu Real-Time Hot Searches"},
	"/v2/baidu/teleplay":             {Emoji: "📺", Zh: "百度电视剧榜", En: "Baidu TV Drama Rankings"},
	"/v2/baidu/tieba":                {Emoji: "💬", Zh: "百度贴吧话题榜", En: "Baidu Tieba Topic Rankings"},
	"/v2/weather/realtime":           {Emoji: "🌤️", Zh: "实时天气", En: "Real-Time Weather"},
	"/v2/weather/forecast":           {Emoji: "☔", Zh: "天气预报", En: "Weather Forecast"},
	"/v2/ncm-rank/list":              {Emoji: "🎧", Zh: "网易云榜单列表", En: "NetEase Cloud Music Rankings"},
	"/v2/ncm-rank/:id":               {Emoji: "🎶", Zh: "网易云榜单详情", En: "NetEase Cloud Music Ranking Details"},
	"/v2/color/random":               {Emoji: "🎨", Zh: "随机颜色/颜色转换", En: "Random Color and Color Conversion"},
	"/v2/color/palette":              {Emoji: "🌈", Zh: "配色方案/色彩搭配", En: "Color Palette and Matching"},
	"/v2/lyric":                      {Emoji: "🎼", Zh: "歌词搜索", En: "Lyric Search"},
	"/v2/fuel-price":                 {Emoji: "⛽", Zh: "汽油价格", En: "Fuel Prices"},
	"/v2/gold-price":                 {Emoji: "🥇", Zh: "黄金价格", En: "Gold Prices"},
	"/v2/olympics":                   {Emoji: "🏅", Zh: "奥运奖牌榜", En: "Olympic Medal Table"},
	"/v2/olympics/events":            {Emoji: "🏟️", Zh: "", En: ""},
	"/v2/douban/weekly/movie":        {Emoji: "🎞️", Zh: "豆瓣全球口碑电影榜", En: "Douban Weekly Global Movie Picks"},
	"/v2/douban/weekly/tv_chinese":   {Emoji: "📺", Zh: "豆瓣华语口碑剧集榜", En: "Douban Weekly Chinese TV Picks"},
	"/v2/douban/weekly/tv_global":    {Emoji: "🌍", Zh: "豆瓣全球口碑剧集榜", En: "Douban Weekly Global TV Picks"},
	"/v2/douban/weekly/show_chinese": {Emoji: "🎭", Zh: "豆瓣华语口碑综艺榜", En: "Douban Weekly Chinese Variety Picks"},
	"/v2/douban/weekly/show_global":  {Emoji: "🌐", Zh: "豆瓣全球口碑综艺榜", En: "Douban Weekly Global Variety Picks"},
	"/v2/og":                         {Emoji: "🔗", Zh: "链接 OG 信息", En: "Link Open Graph Info"},
	"/v2/hash":                       {Emoji: "#️⃣", Zh: "哈希/解压/压缩", En: "Hash, Decompression, and Compression"},
	"/v2/fanyi":                      {Emoji: "🌏", Zh: "在线翻译（支持 109 种语言）", En: "Online Translation (109 Languages)"},
	"/v2/fanyi/langs":                {Emoji: "🈯", Zh: "在线翻译支持的语言列表", En: "Supported Translation Languages"},
	"/v2/beta/kuan":                  {Emoji: "🧪", Zh: "", En: ""},
	"/v2/beta/qq/profile":            {Emoji: "👤", Zh: "", En: ""},
	"/v2/exchange_rate":              {Emoji: "💱", Zh: "", En: ""},
	"/v2/today_in_history":           {Emoji: "📅", Zh: "", En: ""},
	"/v2/maoyan":                     {Emoji: "🎬", Zh: "", En: ""},
	"/v2/baidu/realtime":             {Emoji: "🔥", Zh: "", En: ""},
	"/v2/weather":                    {Emoji: "🌦️", Zh: "", En: ""},
	"/v2/ncm-rank":                   {Emoji: "🎧", Zh: "", En: ""},
	"/v2/color":                      {Emoji: "🎨", Zh: "", En: ""},
}

// LocalizeEndpoints keeps only supported endpoints and replaces names with Chinese menu labels.
func LocalizeEndpoints(endpoints []Endpoint) []Endpoint {
	return LocalizeEndpointsForLanguage(endpoints, "zh")
}

// LocalizeEndpointsForLanguage keeps only supported endpoints and replaces names with localized menu labels.
func LocalizeEndpointsForLanguage(endpoints []Endpoint, language string) []Endpoint {
	localized := make([]Endpoint, 0, len(endpoints))
	seen := make(map[string]struct{}, len(endpoints))
	for _, ep := range endpoints {
		path := strings.TrimSpace(ep.Path)
		if path == "" {
			continue
		}

		menuText, ok := EndpointMenuTexts[lookupEndpointPath(path)]
		name := localizedEndpointName(menuText, language)
		if !ok || name == "" {
			continue
		}
		if menuText.Emoji != "" {
			name = strings.TrimSpace(menuText.Emoji + " " + name)
		}
		if _, exists := seen[path]; exists {
			continue
		}
		seen[path] = struct{}{}

		localized = append(localized, Endpoint{Name: name, Path: path})
	}
	return localized
}

func localizedEndpointName(menuText EndpointMenuText, language string) string {
	switch language {
	case "en":
		if name := strings.TrimSpace(menuText.En); name != "" {
			return name
		}
	}
	return strings.TrimSpace(menuText.Zh)
}

func lookupEndpointPath(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if u, err := url.Parse(endpoint); err == nil && u.Scheme != "" && u.Host != "" {
		return u.Path
	}
	if i := strings.IndexAny(endpoint, "?#"); i >= 0 {
		return endpoint[:i]
	}
	return endpoint
}
