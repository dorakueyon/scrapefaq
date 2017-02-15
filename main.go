package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Crawler struct {
	models []Model
	url    string
	domain string
}

type Model struct {
	Number    int
	ModelName string
	ModelSku  string
	ModelFaq  string
	ModelUrl  string
}

func (m *Model) GetRow() (row []string) {
	row = make([]string, 5)
	row[0] = strconv.Itoa(m.Number)
	row[1] = m.ModelName
	row[2] = m.ModelSku
	row[3] = m.ModelUrl
	row[4] = m.ModelFaq
	return
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) SetUrl(country string, environment string, language string) (err error) {
	var url string
	if environment != "staging" && environment != "live" {
		fmt.Printf("'-e'オプションのあとは、'staging'か'live'")
		os.Exit(0)
	}
	if environment == "live" {
		url = "http://www.oki.com/"
	} else {
		url = "http://10.253.246.78/"
	}
	c.domain = url
	url = url + country + "/printing/"
	if language != "" {
		url = url + language + "/"
	}
	c.url = url
	return
}

func (c *Crawler) NewCrawler() (err error) {
	faq_url := c.url + "support/faq/"
	doc, err := goquery.NewDocument(faq_url)
	if err != nil {
		panic(err)
	}

	var urls []string
	doc.Find(".column4").Find("a").Each(func(_ int, s *goquery.Selection) {
		html, _ := s.Attr("href")
		doc2, err := goquery.NewDocument(c.domain + html)
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
		if (i+1)%10 != 0 {
			fmt.Printf("-")
		} else {
			fmt.Printf("%d", i+1)
		}
	}
	close(resultCh)

	return
}

func (c *Crawler) crawl(url string, i int, resultCh chan []Model) {
	doc, err := goquery.NewDocument(c.domain + url)
	if err != nil {
		panic(err)
	}
	var models []Model
	var model Model

	model.Number = i + 1
	model.ModelUrl = c.domain + url
	model.ModelName = doc.Find("h1").Text()
	// extract SKU from URL
	splitedUrl := strings.Split(url, "/")
	model.ModelSku = splitedUrl[len(splitedUrl)-2]
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

func (c *Crawler) StoreCSV(path string) (err error) {
	c.SortModel()
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = file.Truncate(0)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)
	writer.UseCRLF = true
	// add first row of csv
	writer.Write([]string{"#", "MODEL_NAME", "MODEL_SKU", "MODEL_URL", "FAQ"})
	for _, v := range c.models {
		writer.Write(v.GetRow())
	}

	writer.Flush()
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
	for _, v := range models[1:] {
		if v.Number > pivot.Number {
			right = append(right, v)
		} else {
			left = append(left, v)
		}
	}
	left = SortData(left)
	right = SortData(right)
	ret = append(left, pivot)
	ret = append(ret, right...)
	return ret
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		country     string
		environment string
		language    string
	)
	flag.StringVar(&country, "c", "eu", "処理する国。(e.g.'eu')")
	flag.StringVar(&environment, "e", "live", "'staging'か、'live'を選択")
	flag.StringVar(&language, "l", "", "複数言語サイトは、言語コードを指定")
	flag.Parse()
	c := NewCrawler()
	err := c.SetUrl(country, environment, language)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.url)
	err = c.NewCrawler()
	if err != nil {
		panic(err)
	}
	var save_name string
	if language != "" {
		save_name = "./" + country + "_" + language + "_" + environment + "_faqId.csv"
	} else {
		save_name = "./" + country + "_" + environment + "_faqId.csv"
	}
	err = c.StoreCSV(save_name)
	if err != nil {
		panic(err)
	}
}
