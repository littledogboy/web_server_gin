package spiders

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"strings"
)

type MenuData struct {
	Sections []ItemSection `json:"sections"`
}

type ItemSection struct {
	Title string `json:"title"`
	Items []Item `json:"items"`
}

func MenuItems(callback func(data MenuData, err error)) {
	menuData := MenuData{Sections: []ItemSection{}}
	var e error

	MRTMenuItems(func(items ItemSection, err error) {
		menuData.Sections = append(menuData.Sections, items)
		e = errors.Join(err)
	})

	callback(menuData, e)
}

func MRTMenuItems(callback func(items ItemSection, err error)) {
	sectionItems := ItemSection{Title: Meirentu.Name, Items: []Item{}}

	c := colly.NewCollector()

	selector := "body > div:nth-child(1) > nav > ul > li > ul > li > a"
	c.OnHTML(selector, func(h *colly.HTMLElement) {
		href := h.Attr("href")

		if !strings.HasPrefix(href, "https") {
			href = getDomanFromElement(h) + href
		}
		title := h.Text

		item := Item{Href: href, Title: title}
		sectionItems.Items = append(sectionItems.Items, item)
	})

	c.OnScraped(func(r *colly.Response) {
		callback(sectionItems, nil)
	})

	c.OnError(func(r *colly.Response, err error) {
		callback(sectionItems, err)
	})

	c.Visit(Meirentu.Doman)
}
