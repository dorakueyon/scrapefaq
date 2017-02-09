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

		}
	})

	// Starting goroutine
	resultCh := make(chan []Model)
	for _, url := range urls {
		go c.crawl(url, resultCh)
	}

	for i := 1; i < pageNum+1; i++ {
		gs := <-resultCh
		c.games = append(c.games, gs...)
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

func (c *Crawler) extractDiscount(discount string) (int, error) {
	re, err := regexp.Compile("[0-9]+")
	if err != nil {
		panic(err)
	}
	exDiscount := re.FindString(discount)
	return strconv.Atoi(exDiscount)
}

func (c *Crawler) storeCSV(path string) (err error) {
	c.SortGames()
	return
}
func (c *Crawler) SortGames() {
	c.games = sortData(c.games)
}

func sortData(games []Game) (ret []Game) {
	if len(games) == 0 {
		return games
	}
	pivot := games[0]

	var left []Game
	var right []Game

	for _, v := range games[1:] {
		if v.Number > pivot.Number {
			right = append(right, v)
		} else {
			left = append(left, v)
		}
	}
	left = sortData(left)
	right = sortData(right)
	ret = append(left, pivot)
	ret = append(ret, right...)
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := NewCrawler()
	err := c.StartCrawl()
	if err != nil {
		panic(err)
	}
