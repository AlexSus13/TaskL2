package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
)

func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		//title := s.Find("a").Text()
		//fmt.Printf("Review %d: %s\n", i, title)
		link, err := s.Attr("href")
		fmt.Println("err = ", err)
		fmt.Printf("link %s\n", link)
	})
}

func main() {
	ExampleScrape()
}
