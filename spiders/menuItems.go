package spiders

import (
	"github.com/gocolly/colly/v2"
	"strings"
)

type Items struct {
	Items []Item `json:"items"`
}

func MRTMenuItems(callback func(items Items, err error)) {
	items := Items{Items: []Item{}}

	c := colly.NewCollector()

	selector := "body > div:nth-child(1) > nav > ul > li > ul > li > a"
	c.OnHTML(selector, func(h *colly.HTMLElement) {
		href := h.Attr("href")

		if !strings.HasPrefix(href, "https") {
			href = getDomanFromElement(h) + href
		}
		title := h.Text

		item := Item{Href: href, Title: title}
		items.Items = append(items.Items, item)
	})

	c.OnScraped(func(r *colly.Response) {
		callback(items, nil)
	})

	c.OnError(func(r *colly.Response, err error) {
		callback(items, err)
	})

	c.Visit(Meirentu.Doman)
}
