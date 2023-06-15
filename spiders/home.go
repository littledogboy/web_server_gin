package spiders

import (
	// "bytes"
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"

	"net/url"
	"strconv"
	"strings"
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

func MRTHomeSpider(desUrl string, page string, callback func(Home, error)) {
	var newUrl string = desUrl
	if page != "" {
		newUrl = Meirentu.Doman + "index/" + page + ".html"
	}
	MRTDesURLSpider(newUrl, page, Meirentu.Refer, Meirentu.ReferValue, Meirentu_Home_Selector, callback)
}

func GroupSpider(href string, page string, callback func(Home, error)) {
	if strings.Contains(href, Meirentu.Doman) || strings.Contains(href, Fulitu.Doman) {
		MRTGroupSpider(href, page, callback)
	} else if strings.Contains(href, Bestprettygirl.Doman) {
		BPGGroupSpider(href, page, callback)
	}
}

func MRTGroupSpider(href string, page string, callback func(Home, error)) {
	var newHref string
	if strings.Contains(href, Meirentu.Doman) {
		array1 := strings.Split(href, "-")
		newHref = array1[0] + "-" + page + ".html"
		MRTDesURLSpider(newHref, page, Meirentu.Refer, Meirentu.ReferValue, Meirentu_Group_Selector, callback)
	} else if strings.Contains(href, Fulitu.Doman) {
		if strings.Contains(href, "-") {
			array1 := strings.Split(href, "-")
			newHref = array1[0] + "-" + page + ".html"
			MRTDesURLSpider(newHref, page, "", "", Meirentu_Group_Selector, callback)
		} else {
			newHref = strings.Replace(href, ".html", "-"+page+".html", 1)
			MRTDesURLSpider(newHref, page, "", "", Meirentu_Group_Selector, callback)
		}
	}
}

func BPGGroupSpider(href string, page string, callback func(Home, error)) {
	var newHref = href

	if page != "" {
		newHref = href + "page" + "/" + page
	}
	BPGDesURLSpider(newHref, page, "", "", Bestprettygirl_Thumbnail_Selector, callback)
}

func MRTTagPageSpider(href string, page string, callback func(Home, error)) {
	urlStr, _ := url.JoinPath(Meirentu.Doman, href)
	url, _ := url.PathUnescape(urlStr)

	intPage, err := strconv.Atoi(page)
	if err == nil && intPage > 1 {
		home := Home{}
		callback(home, nil)
	} else {
		MRTDesURLSpider(url, page, Meirentu.Refer, Meirentu.ReferValue, Meirentu_TagPage_Selector, callback)
	}
}

func MRTDesURLSpider(desUrl string, page string, refer string, value string, selector string, callback func(Home, error)) {
	home := Home{
		Recommends: []Item{},
	}

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

	// 注册请求回调
	c.OnRequest(func(r *colly.Request) {
		// log.Println(r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		// log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
		// log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	})

	// 注册 html 回调
	c.OnHTML(selector, func(li *colly.HTMLElement) {
		href := li.ChildAttr("a", "href")
		if !strings.HasPrefix(href, "https") {
			doman := getDomanFromElement(li)
			href = doman + href
		}
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
		log.Println("请求的URL：", r.Request.URL, "失败的响应：", r, "\n错误：", err)
		callback(home, err)
	})

	c.Visit(desUrl)
}

func BPGDesURLSpider(desUrl string, page string, refer string, value string, selector string, callback func(Home, error)) {
	home := Home{
		Recommends: []Item{},
	}

	// 创建采集器
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	// 注册请求回调
	c.OnRequest(func(r *colly.Request) {
	})

	// 注册 html 回调
	c.OnHTML(selector, func(h *colly.HTMLElement) {
		href := h.Attr("href")
		newHref, _ := url.PathUnescape(href)
		img := h.ChildAttr("div > img", "src")
		title := h.ChildAttr("div > img", "alt")

		item := Item{
			Href:  newHref,
			Img:   img,
			Title: title,
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
		callback(home, err)
	})

	c.Visit(desUrl)
}

func getDomanFromElement(h *colly.HTMLElement) string {
	scheme := h.Request.URL.Scheme
	host := h.Request.URL.Host
	return scheme + "://" + host
}
