package scripts

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const protocolIT = "protocolIT"

var protocolCSV = fmt.Sprintf("files/%s/%s.csv", protocolIT, protocolIT)

var titleIDName = "#offer-title"

func ScrapProtocol(path string, done chan struct{}) {
	protocolListings := []ListingData{}

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

	doc.Find(titleIDName).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s - found position: %s\n", protocolIT, s.Text())
		protocolListings = append(protocolListings, ListingData{
			Title: s.Text(),
		})
	})

	csvFile, err := os.Create(protocolCSV)
	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}

	if err := WriteListingsToCSV(protocolListings, csvFile); err != nil {
		fmt.Println("error during appending data to csv file: ", err)
	}

	defer csvFile.Close()
	done <- struct{}{}
}
