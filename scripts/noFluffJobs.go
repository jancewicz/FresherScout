package scripts

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const noFluffJobs = "noFluffJobs"

var nfjCSV = fmt.Sprintf("files/%s/%s.csv", noFluffJobs, noFluffJobs)

var titleClassName = ".posting-title__position.ng-star-inserted" // CSS class for job name

// Function scraps existing in 'files' directory HTML for NFJ platform. Then Saves results in CSV file.
func ScrapNFJHTML(path string, done chan struct{}) {
	var nfjListings = []ListingData{}

	htmlFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open html file: %v", err)
	}
	defer htmlFile.Close()

	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		fmt.Println("Cannot read HTML file: ", err)
		return
	}

	doc.Find(titleClassName).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s - found position: %s\n", noFluffJobs, s.Text())
		nfjListings = append(nfjListings, ListingData{
			Title: s.Text(),
		})
	})

	csvFile, err := os.Create(nfjCSV)
	if err != nil {
		fmt.Println("Cannot create file", err)
	}
	defer csvFile.Close()

	if err := WriteListingsToCSV(nfjListings, csvFile); err != nil {
		fmt.Println("error during appending data to csv file: ", err)
	}

	done <- struct{}{}
}
