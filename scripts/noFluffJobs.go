package scripts

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type NoFluffJobsData struct {
	Position string
}

var NFJPathCSV = "files/noFluffJobs/nofluffjobs.csv"

var positionClass = ".posting-title__position.ng-star-inserted" // CSS class for job name

// Function scraps existing in 'files' directory HTML for NFJ platform. Then Saves results in CSV file.
func ScrapNFJHTML(path string, done chan struct{}) {
	var offers = []NoFluffJobsData{}

	htmlFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open html file: %v", err)
	}
	defer htmlFile.Close()

	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		log.Fatalf("Cannot read HTML file: %v\n", err)
	}

	doc.Find(positionClass).Each(func(i int, s *goquery.Selection) {
		fmt.Println("found position: ", s.Text())
		offers = append(offers, NoFluffJobsData{
			Position: s.Text(),
		})
	})

	csvFile, err := os.Create(NFJPathCSV)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	headers := []string{
		"Position",
	}

	err = csvwriter.Write(headers)
	if err != nil {
		log.Fatalln("Error writing record to CSV:", err)
	}

	for _, offer := range offers {
		row := []string{offer.Position}
		if err := csvwriter.Write(row); err != nil {
			log.Fatalln("Error writing record to CSV:", err)
		}
	}
	done <- struct{}{}
}
