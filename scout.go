package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Article struct {
	Title,
	Content string
}

func main() {
	fmt.Println("Lets GO scout!")

	url := "https://www.scrapethissite.com/pages/"
	var articles []Article

	c := colly.NewCollector(
		colly.AllowedDomains(url),
	)

	c.OnHTML("div.page", func(e *colly.HTMLElement) {
		article := Article{}

		article.Title = e.ChildText(".page-title")
		article.Content = e.ChildText(".lead session-desc")

		articles = append(articles, article)
	})

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("articles.csv")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		headers := []string{
			"Title",
			"Content",
		}
		writer.Write(headers)

		for _, article := range articles {
			articleSlice := []string{
				article.Title,
				article.Content,
			}

			writer.Write(articleSlice)
		}
		defer writer.Flush()
	})

	c.Visit(url)
}
