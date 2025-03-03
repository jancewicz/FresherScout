package scripts

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const bulldogJobs = "bulldogJobs"

var bulldogCSV = fmt.Sprintf("files/%s/%s.csv", bulldogJobs, bulldogJobs)

var bulldogJobsTitleClass = ".JobListItem_item__title__278xz"

func ScrapBulldogJobs(path string) {
	bulldogListings := []ListingData{}

	htmlFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open the file: %v", err)
	}
	defer htmlFile.Close()

	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		fmt.Printf("Cannot read HTML file: %v\n", err)
		return
	}

	doc.Find(bulldogJobsTitleClass).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s - found position: %s\n", bulldogJobs, s.Text())
		bulldogListings = append(bulldogListings, ListingData{
			Title: s.Text(),
		})
	})

	csvFile, err := os.Create(bulldogCSV)
	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}
	defer csvFile.Close()

	if err := WriteListingsToCSV(bulldogListings, csvFile); err != nil {
		fmt.Println("error during appending data to csv file: ", err)
	}

	defer csvFile.Close()
}
