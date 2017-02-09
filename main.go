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
	if err != nil {
		return
	}

	//Getting the number of pages.
	var pageNum int
	doc.Find(".search_pagination_right").Children().Each(func(i int, s *goquery.Selection) {
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

	for i := 1; i < pageNum+1; i++ {
		gs := <-resultCh
		c.games = append(c.games, gs...)
	}
	close(resultCh)

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
		panic(err)
	}
	// Scraping
	var games []Game
	doc.Find(".search_result_row").Each(func(_ int, s *goquery.Selection) {
		var game Game

		game.Name = s.Find(".title").Text()
		game.ReleaseDate = s.Find(".search_released").Text()

		// Getting discount rate.
		game.DiscountRate, _ = c.extractDiscount(s.Find(".search_discount").Find("span").Text())

		game.Number = elementNum
		elementNum++

		games = append(games, game)
	})
	resultCh <- games
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
	err2 := c.storeCSV("data.csv")
	if err2 != nil {
		panic(err2)
	}
}
