package spiders

import (
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"

	"github.com/gocolly/colly/v2/proxy"
	"log"
)

type TagData struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Title    string `json:"title"`
	Color    string `json:"color"`
	FontSize int    `json:"fontSize"`
	Href     string `json:"href"`
}

func MRTTagsSpider(callback func(tagData TagData, err error)) {
	tagData := TagData{Tags: []Tag{}}

	// 创建采集器
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("meirentu.cc", "fulitu.me"),
	)
	extensions.RandomUserAgent(c)

	// proxies
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:7890")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	selector := Meirentu_Tags_Selector
	c.OnHTML(selector, func(h *colly.HTMLElement) {
		class := h.Attr("class")
		splitA := strings.Split(class, " ")
		fontSize := TagFontSizeMap[splitA[0]]
		fontColor := TagFontColorMap[splitA[1]]
		title := h.Text
		href := h.Attr("href")

		tag := Tag{
			Title:    title,
			Color:    fontColor,
			FontSize: fontSize,
			Href:     href,
		}

		tagData.Tags = append(tagData.Tags, tag)
	})

	c.OnScraped(func(r *colly.Response) {
		callback(tagData, nil)
	})

	c.OnError(func(r *colly.Response, err error) {
		callback(tagData, err)
	})

	c.Visit(Meirentu_TagsPage)
}
