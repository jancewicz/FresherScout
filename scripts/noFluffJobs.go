package scripts

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type NoFluffJobsData struct {
	Position string
}

var url = "https://nofluffjobs.com/pl/Golang"
var positionClass = "a:visited.posting-title__position"
var noFluffDomain = "nofluffjobs.com"

func ScrapNoFluffJobs() {
	var noFluffJobsOffers []NoFluffJobsData

	c := colly.NewCollector(
		colly.AllowedDomains(noFluffDomain),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"

	c.OnHTML(positionClass, func(e *colly.HTMLElement) {
		data := NoFluffJobsData{}
		data.Position = e.Text

		noFluffJobsOffers = append(noFluffJobsOffers, data)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Visited: %s, status code: %d ", r.Request.URL, r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("nofluffjobs.csv")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}

		defer file.Close()

		writer := csv.NewWriter(file)
		headers := []string{
			"Position",
		}
		writer.Write(headers)

		for _, offer := range noFluffJobsOffers {
			offerSlice := []string{
				offer.Position,
			}

			writer.Write(offerSlice)
		}
		defer writer.Flush()
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal("Cannot visit No Fluff Jobs", err)
	}
}
