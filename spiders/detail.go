package spiders

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gocolly/colly/v2"
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

func DetailViewSpider(urlString string, callback func(AlbumDetail)) {

	c1 := colly.NewCollector(
		colly.AllowedDomains("meirentu.cc"),
		colly.Async(true),
	)

	err := c1.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 100 * time.Millisecond,
		Parallelism: 30,
	})

	imageCollector := c1.Clone()

	if err != nil {
		println("添加限制出错: ", err)
	}

	c1.OnRequest(func(r *colly.Request) {
		r.Headers.Add("referer", "https://meirentu.cc/")
		// println("正在访问:", r.URL.String())
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
		// println("c1 抓取结束 \n--------------------\n")
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("请求失败的URL：", r.Request.URL, "失败的响应：", r, "\n错误：", err)
	})

	//////////////////
	albumDetail := AlbumDetail{
		Info:   "",
		Images: []Image{},
	}

	imageCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Add("referer", "https://meirentu.cc/")
		// println("imageCollector 正在访问:", r.URL.String())
	})

	// // 获取图片
	imgQuery := "body > div.main > div > div > div:nth-child(3) > div > div > img"
	imageCollector.OnHTML(imgQuery, func(img *colly.HTMLElement) {
		imageSrc := img.Attr("src")
		if imageSrc != "" {
			// 获取图片
			getImageSize(imageSrc, func(width, height int) {
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
		// println("imageCollector 抓取结束", r.Request.URL.String(), "\n--------------------\n")
	})

	// 访问
	c1.Visit(urlString)
	c1.Wait()
	imageCollector.Wait()

	// 排序
	sort.Slice(albumDetail.Images, func(i, j int) bool {
		return albumDetail.Images[i].Src < albumDetail.Images[j].Src
	})
	callback(albumDetail)
	// println("详情页抓取结束\n")
	// println("------------\n")
}

func getImageSize(imageStr string, callback func(int, int)) {

	referer := "referer"
	doman := "https://meirentu.cc"

	req, _ := http.NewRequest("GET", imageStr, nil)
	req.Header.Set(referer, doman)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	m, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	g := m.Bounds()

	// Get height and width
	height := g.Dy()
	width := g.Dx()

	callback(width, height)
}