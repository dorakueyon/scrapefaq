package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type QueryResult struct {
	Entries []*entry
}

func main() {
	// Open FAQ category page
	url := "www-odc-ori.oki.com/jp/printing/support/faq/color/index.html"
	// Get FAQ detail pages
	result, err := crawl(url)
	// open each faq detail page
	// get faq id
	// store faq id to result
}

func crawl(url string) (QueryResult, error) {
	doc, err := goquery.NewDocument(url)
	if err != {
		QueryResult{
		ProductName string
		ProductURL  string
		}
	}

}
