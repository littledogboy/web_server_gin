package spiders

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

type TagData struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Title    string `json:"title"`
	Color    string `json:"color"`
	FontSize int    `json:"fontSize"`
	Href     string `json:"Href"`
}

func MRTTagsSpider(callback func(tagData TagData, err error)) {
	tagData := TagData{Tags: []Tag{}}

	c := colly.NewCollector()

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
