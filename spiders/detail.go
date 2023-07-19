package spiders

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
)

type AlbumDetail struct {
	Info   string  `json:"info"`
	Images []Image `json:"images"`
}

type Image struct {
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func DetailViewSpider(urlString string, callback func(AlbumDetail, error)) {
	site := WebsiteFromURLString(urlString)

	switch site {
	case MeirentuSite, FulituSite:
		MRTDetailViewSpider(urlString, callback)
	case BestprettygirlSite:
		BPGDetailViewSpider(urlString, callback)
	}
}

func BPGDetailViewSpider(urlString string, callback func(AlbumDetail, error)) {
	var e error

	newURLString, _ := url.PathUnescape(urlString)
	albumDetail := AlbumDetail{
		Info:   "",
		Images: []Image{},
	}

	c := colly.NewCollector(
		colly.Async(true),
	)
	extensions.RandomUserAgent(c)

	c.Limit(&colly.LimitRule{
		DomainGlob:  Bestprettygirl.Doman,
		RandomDelay: 100 * time.Millisecond,
		Parallelism: 30,
	})

	imageCollector := c.Clone()

	c.OnHTML(BestPrettyGirl_Detail_Selector, func(h *colly.HTMLElement) {
		imageSrc := h.Attr("src")
		if imageSrc != "" && strings.HasSuffix(imageSrc, "jpg") {
			imageCollector.Visit(imageSrc)
		}
	})

	c.OnScraped(func(r *colly.Response) {
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("请求失败的URL：", r.Request.URL)
		e = errors.Join(err)
	})

	imageCollector.OnResponse(func(r *colly.Response) {
		m, _, err := image.Decode(bytes.NewReader(r.Body))
		if err != nil {
			log.Println(err)
			e = errors.Join(err)
			return
		}

		g := m.Bounds()

		src := r.Request.URL.String()
		newSrc, _ := url.PathUnescape(src)
		// Get height and width
		height := g.Dy()
		width := g.Dx()
		image := Image{
			Src:    newSrc,
			Width:  int(width),
			Height: int(height),
		}
		albumDetail.Images = append(albumDetail.Images, image)
	})

	imageCollector.OnError(func(r *colly.Response, err error) {
		log.Println("请求失败的URL：", r.Request.URL)
		e = errors.Join(err)
	})

	c.Visit(newURLString)
	c.Wait()
	imageCollector.Wait()

	sort.Slice(albumDetail.Images, func(i, j int) bool {
		return albumDetail.Images[i].Src < albumDetail.Images[j].Src
	})

	callback(albumDetail, e)
}

func MRTDetailViewSpider(urlString string, callback func(AlbumDetail, error)) {
	refer, referValue := getReferValueFromURLString(urlString)

	albumDetail := AlbumDetail{
		Info:   "",
		Images: []Image{},
	}

	// 创建采集器
	c1 := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("meirentu.cc", "fulitu.me"),
		colly.Async(true),
	)
	extensions.RandomUserAgent(c1)

	// proxies
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:7891")
	if err != nil {
		log.Fatal(err)
	}
	c1.SetProxyFunc(rp)

	c1.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 100 * time.Millisecond,
		Parallelism: 30,
	})

	imageCollector := c1.Clone()

	c1.OnRequest(func(r *colly.Request) {
		if refer != "" && referValue != "" {
			r.Headers.Add(refer, referValue)
		}
	})

	// 解析 page 获取所有分页链接
	pageQuery := "div.content_left > div.page"
	c1.OnHTML(pageQuery, func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, h *colly.HTMLElement) {
			href := h.Attr("href")
			text := h.Text
			if text != "下页" {
				fullLink := h.Request.AbsoluteURL(href)
				imageCollector.Visit(fullLink)
			}
		})
	})

	c1.OnScraped(func(r *colly.Response) {
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("请求失败的URL：", r.Request.URL, "失败的响应：", r, "\n错误：", err)
		callback(albumDetail, err)
	})

	imageCollector.OnRequest(func(r *colly.Request) {
		if refer != "" && referValue != "" {
			r.Headers.Add(refer, referValue)
		}
	})

	// // 获取图片
	imgQuery := "body > div.main > div > div > div:nth-child(3) > div > div > img"
	imageCollector.OnHTML(imgQuery, func(img *colly.HTMLElement) {
		imageSrc := img.Attr("src")
		if imageSrc != "" {
			// 获取图片
			getImageSize(imageSrc, refer, referValue, func(width, height int) {
				image := Image{
					Src:    imageSrc,
					Width:  width,
					Height: height,
				}
				albumDetail.Images = append(albumDetail.Images, image)
			})
		}
	})

	imageCollector.OnScraped(func(r *colly.Response) {
	})

	imageCollector.OnError(func(r *colly.Response, err error) {
		callback(albumDetail, err)
	})

	// 访问
	c1.Visit(urlString)
	c1.Wait()
	imageCollector.Wait()

	// 排序
	sort.Slice(albumDetail.Images, func(i, j int) bool {
		return albumDetail.Images[i].Src < albumDetail.Images[j].Src
	})

	callback(albumDetail, nil)
}

func getImageSize(imageStr string, refer string, value string, callback func(int, int)) {
	req, _ := http.NewRequest("GET", imageStr, nil)
	if refer != "" && value != "" {
		req.Header.Set(refer, value)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	m, _, err := image.Decode(resp.Body)
	if err != nil {
		return
	}
	g := m.Bounds()

	// Get height and width
	height := g.Dy()
	width := g.Dx()

	callback(width, height)
}

func getReferValueFromURLString(urlStr string) (refer string, value string) {
	var r, v string
	u, _ := url.Parse(urlStr)
	hostname := u.Hostname()
	if strings.Contains(Meirentu.Doman, hostname) {
		r = Meirentu.Refer
		v = Meirentu.ReferValue
	}
	return r, v
}
