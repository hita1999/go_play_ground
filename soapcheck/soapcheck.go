package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	//targetURL
	url := "https://www.uec-programming.com/e-learning/irohaboard/"

	c := colly.NewCollector()

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.Visit(url)
}
