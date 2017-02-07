package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"sync"
)

type Result struct {
	Title string
	Url   string
}

func GetPages(url string) []Result {
	results := []Result{}
	result := Result{url, url}
	results = append(results, result)
	return results
}

func GoGet(urls []string) <-chan []Result {
	var wg sync.WaitGroup
	ch := make(chan []Result)
	go func() {
		fmt.Printf("test")
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				//ch <- GetPages(url)
				fmt.Println(url)
				ch <- GetPages(url)
				wg.Done()
			}(url)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

func main() {
	// Open FAQ category page
	url := "http://www-odc-ori.oki.com/jp/printing/support/faq/index.html"
	// Get FAQ detail pages
	urls := crawl(url)
	ch := GoGet(urls)
	// open each faq detail page
	// get faq id
	// store faq id to result
	fmt.Println(ch)
}

func crawl(url string) []string {
	results := []string{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".column4").Find("a").Each(func(_ int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			results = append(results, url)
		}
	})
	return results
}
