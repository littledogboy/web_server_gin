package spiders

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
	"log"
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
		if err == nil {
			menuData.Sections = append(menuData.Sections, items)
		}
		e = errors.Join(err)
	})

	BestPrettyGirlMenuItems(func(items ItemSection, err error) {
		if err == nil {
			menuData.Sections = append(menuData.Sections, items)
		}
		e = errors.Join(err)
	})

	callback(menuData, e)
}

func MRTMenuItems(callback func(items ItemSection, err error)) {
	sectionItems := ItemSection{Title: Meirentu.Name, Items: []Item{}}

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

func BestPrettyGirlMenuItems(callback func(items ItemSection, err error)) {
	sectionItems := ItemSection{Title: Bestprettygirl.Name, Items: []Item{}}

	c := colly.NewCollector()

	// proxies
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:7890")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	extensions.RandomUserAgent(c)

	selector := Bestprettygirl_Menu_Selector
	c.OnHTML(selector, func(h *colly.HTMLElement) {
		href := h.Attr("href")
		title := h.Text

		item := Item{Href: href, Title: title}

		if title != "Video" {
			sectionItems.Items = append(sectionItems.Items, item)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		callback(sectionItems, nil)
	})

	c.OnError(func(r *colly.Response, err error) {
		callback(sectionItems, err)
	})

	c.Visit(Bestprettygirl.Doman)
}
