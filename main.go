package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"runtime"
	"strings"
)

type Crawler struct {
	models []Model
}

type Model struct {
	Number    int
	ModelName string
	ModelSku  string
	ModelFaq  string
	ModelUrl  string
}

const URL = "http://www.oki.com/jp/printing/support/faq/index.html"
const DOMAIN = "http://www.oki.com/"

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) NewCrawler() (err error) {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		panic(err)
	}

	var urls []string
	doc.Find(".column4").Find("a").Each(func(_ int, s *goquery.Selection) {
		html, _ := s.Attr("href")
		doc2, err := goquery.NewDocument(DOMAIN + html)
		if err != nil {
			panic(err)
		}
		doc2.Find(".col-xs-12" + ".col-sm-9").Find("a").Each(func(_ int, s *goquery.Selection) {
			url, ok := s.Attr("href")
			if ok {
				urls = append(urls, url)
			}
		})
	})

	// start crawl
	fmt.Println("Start crawling:", len(urls), "models")
	resultCh := make(chan []Model)
	for i, url := range urls {
		go c.crawl(url, i, resultCh)
	}

	for i := 0; i < len(urls); i++ {
		gs := <-resultCh
		c.models = append(c.models, gs...)
	}
	close(resultCh)
	fmt.Println(c.models)

	return
}

func (c *Crawler) crawl(url string, i int, resultCh chan []Model) {
	doc, err := goquery.NewDocument(DOMAIN + url)
	if err != nil {
		panic(err)
	}
	var models []Model
	var model Model

	model.Number = i
	model.ModelUrl = DOMAIN + url
	model.ModelName = doc.Find("h1").Text()
	// extract SKU from URL
	splitedUrl := strings.Split(url, "/")
	model.ModelSku = splitedUrl[len(splitedUrl)-1]
	wrapper := doc.Find(".tabContentsWrapper").Text()
	model.ModelFaq = c.extractFaq(wrapper)
	models = append(models, model)
	resultCh <- models
}

func (c *Crawler) extractFaq(html string) string {
	re, err := regexp.Compile(`p\s\:\s\"(.*)\"`)
	if err != nil {
		panic(err)
	}
	faq := re.FindStringSubmatch(html)
	// if faq cannot be matched return null
	if len(faq) == 1 {
		return ""
	}
	return faq[1]
}

func (c *Crawler) StoreCSV() {
	c.SortModel()
	return
}

func (c *Crawler) SortModel() {
	c.models = SortData(c.models)
}

func SortData(models []Model) (ret []Model) {
	if len(models) == 0 {
		return models
	}
	pivot := models[0]
	var left []Model
	var right []Model
	for _, v := range models {
		if v.Number < pivot.Number {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	SortData(left)
	SortData(right)
	ret = append(left, pivot)
	ret = append(ret, right...)
	return ret
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := NewCrawler()
	err := c.NewCrawler()
	if err != nil {
		panic(err)
	}
	//err2 := c.StoreCSV()
	//if err != nil {
	//	panic(err)
	//}
	c.StoreCSV()
}
