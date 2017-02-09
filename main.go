package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"runtime"
)

type Model struct {
	ModelName string
	ModelSku  string
	Faq       string
	URL       string
}

const URL = "http://www-odc-ori.oki.com/jp/printing/support/faq/color/index.html"

type Crawler struct {
	models []Model
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) StartCrawl() (err error) {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return
	}

	var urls []string

	doc.Find(".col-xs-12" + ".col-sm-9").Find("a").Each(func(_ int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			urls = append(urls, url)
		}
	})

	// Starting goroutine
	resultCh := make(chan []Model)
	for _, url := range urls {
		go c.crawl(url, resultCh)
	}

	for i := 0; i < len(urls); i++ {
		gs := <-resultCh
		c.models = append(c.models, gs...)
	}
	close(resultCh)
	return
}

func (c *Crawler) crawl(url string, resultCh chan []Model) {
	base_url := "http://www.oki.com"
	doc, err := goquery.NewDocument(base_url + url)
	if err != nil {
		panic(err)
	}

	var models []Model

	var model Model
	re, _ := regexp.Compile(`p\s\:\s\"(.*)\"`)
	html := re.FindStringSubmatch(doc)
	model.Faq = html[1]
	model.ModelName = doc.Find("h1").Text()
	model.URL = url
	models = append(models, model)

	resultCh <- models
}

func (c *Crawler) extractFaq(html string) (string, error) {
	re, err := regexp.Compile(`p\s\:\s\"(.*)\"`)

	fmt.Println(html)
	if err != nil {
		panic(err)
	}
	faq := re.FindString(html)
	return faq, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := NewCrawler()
	err := c.StartCrawl()
	if err != nil {
		panic(err)
	}
}
