package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"runtime"
	"strconv"
)

type Crawler struct {
	games []Game
}

type Game struct {
	Number        int    `json:"number"`
	Name          string `json:"name"`
	ReleaseDate   string `json:"releaseDate"`
	DiscountRate  int    `json:"discountRate"`
	NormalPrice   int    `json:"normalPrice"`
	DiscountPrice int    `json:"discountPrice"`
	Rate          int    `json:"rate"`
	Reviewer      int    `json:"reviewer"`
	URL           string `json:"url"`
}

const URL = "http://store.steampowered.com/search/results?sort_by=_ASC&specials=1"

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) StartCrawl() (err error) {
	doc, err := goquery.NewDocument(URL)
	fmt.Println(doc)
	if err != nil {
		return
	}

	//Getting the number of pages.
	var pageNum int
	doc.Find(".search_pagination_right").Children().Each(func(i int, s *goquery.Selection) {
		fmt.Println(s)
		if i == 2 {
			pageNum, err = strconv.Atoi(s.Text())
			if err != nil {
				return
			}
		}
	})

	//Starting goroutine.
	resultCh := make(chan []Game, pageNum)
	for i := 1; i < pageNum+1; i++ {
		url := fmt.Sprintf("%s&page=%d", URL, i)
		go c.crawl(url, resultCh)
	}

	return
}

func (c *Crawler) crawl(url string, resultCh chan []Game) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	//Getting the number of element of a page
	elementNum, err := c.getFirstElementNumber(doc.Find(".search_pagination_left").Text())
	if err != nil {
		fmt.Printf("pani")
		panic(err)
	}
	fmt.Printf("num", elementNum)
}

func (c *Crawler) getFirstElementNumber(paginationLeft string) (int, error) {
	re, err := regexp.Compile(`[0-9]+`)
	if err != nil {
		return 0, err
	}
	pageStr := re.FindString(paginationLeft)

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		panic(err)
	}
	return page, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := NewCrawler()
	err := c.StartCrawl()
	if err != nil {
		panic(err)
	}

}
