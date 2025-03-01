package scripts

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type NoFluffJobsData struct {
	Position string
}

var url = "https://nofluffjobs.com/pl/Golang"
var positionClass = "a:visited.posting-title__position" // CSS class for job name
var noFluffDomain = "nofluffjobs.com"

func checkPositions() bool {
	if _, err := os.Stat("nofluffjobs.csv"); os.IsNotExist(err) {
		fmt.Println("file does not exist")
		return false
	} else {
		file, err := os.Open("nofluffjobs.csv")
		if err != nil {
			log.Fatal("Cannot open file", err)
		}
		defer file.Close()

		csvReader := csv.NewReader(file)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal("Cannot read file", err)
		}

		for _, row := range data {
			for _, col := range row {
				if strings.Contains(col, "Junior") || strings.Contains(col, "junior") || strings.Contains(col, "trainee") || strings.Contains(col, "Trainee") {
					fmt.Println("Junior position found")
					return true
				}
			}
		}
	}
	fmt.Println("Cannot find junior position")
	return false
}

func ScrapNoFluffJobs() {
	done := make(chan bool)
	var noFluffJobsOffers []NoFluffJobsData

	go func() {
		defer close(done)

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
			defer writer.Flush()

			headers := []string{
				"Position",
			}

			if err := writer.Write(headers); err != nil {
				log.Fatal("Cannot write headers to file", err)
			}

			for _, offer := range noFluffJobsOffers {
				offerSlice := []string{
					offer.Position,
				}
				if err := writer.Write(offerSlice); err != nil {
					log.Fatal("Cannot write offer to file", err)
				}
			}
		})

		if err := c.Visit(url); err != nil {
			log.Fatal("Cannot visit No Fluff Jobs", err)
			done <- false
			return
		}

		done <- checkPositions()
	}()
}
