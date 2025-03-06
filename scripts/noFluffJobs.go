package scripts

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const noFluffJobs = "noFluffJobs"

var titleClassName = ".posting-title__position.ng-star-inserted" // CSS class for job name

// Path: path to  html file
func ScrapNFJHTML(path string) []ListingData {
	var listings = []ListingData{}

	doc := CreateGoQuery(path)

	doc.Find(titleClassName).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s - found position: %s\n", noFluffJobs, s.Text())

		listing := ListingData{
			Page:  noFluffJobs,
			Title: s.Text(),
		}

		listings = append(listings, listing)
	})

	return listings
}
