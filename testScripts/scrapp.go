package scrapp

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

func ScrapFirst() {
	url := "https://www.scrapethissite.com/pages/"
	var articles []Article

	// set proxy if needed

	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapethissite.com"),
	)

	c.OnHTML("div.page", func(e *colly.HTMLElement) {
		article := Article{}

		article.Title = e.ChildText(".page-title")
		article.Content = e.ChildText(".lead.session-desc")

		articles = append(articles, article)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Visited: %s, status code: %d ", r.Request.URL, r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("csvs/articles.csv")
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

	err := c.Visit(url)
	if err != nil {
		log.Fatal("Cannot visit site", err)
	}
}

type Product struct {
	Url, Image, Name, Price string
}

func ScrapSecond() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
	)

	var products []Product

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		product := Product{}

		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		products = append(products, product)

	})

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("csvs/products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}
		writer.Write(headers)
		for _, product := range products {
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}
			writer.Write(record)
		}
		defer writer.Flush()
	})

	c.Visit("https://www.scrapingcourse.com/ecommerce")
}
