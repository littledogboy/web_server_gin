package spiders

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Home struct {
	Recommends []Item `json:"recommends"`
}

type Item struct {
	Href  string `json:"href"`
	Img   string `json:"img"`
	Model string `json:"model"`
	Title string `json:"title"`
	Time  string `json:"time"`
}

func RecommendSpider(page string, callback func(Home, error)) {

	// json结构体
	home := Home{
		Recommends: []Item{},
	}

	// 域名
	doman := Meirentu
	// 创建采集器
	c := colly.NewCollector()

	// 注册请求回调
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("referer", doman)
	})

	// 注册 html 回调
	c.OnHTML("body > div.update_area > div > ul > li", func(li *colly.HTMLElement) {
		href := doman + li.ChildAttr("a", "href")
		img := li.ChildAttr("a > img", "src")
		model := li.ChildText("a > div > span")
		title := li.ChildText("div > div.meta-title")
		time := li.ChildText("div > div.meta-post > span:nth-child(2)")

		item := Item{
			Href:  href,
			Img:   img,
			Model: model,
			Title: title,
			Time:  time,
		}

		home.Recommends = append(home.Recommends, item)
	})

	// 抓取结束后回调
	c.OnScraped(func(r *colly.Response) {
		callback(home, nil)
	})

	// 错误回调
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("请求的URL：", r.Request.URL, "失败的响应：", r, "\n错误：", err)
	})

	// https://meirentu.cc/index/1.html
	// 访问
	if page != "" {
		html := doman + "index/" + page + ".html"
		c.Visit(html)
	} else {
		c.Visit(doman)
	}
}
