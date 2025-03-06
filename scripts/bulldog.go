package scripts

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const bulldogJobs = "bulldogJobs"

var bulldogJobsTitleClass = ".JobListItem_item__title__278xz"

// Path: path to  html file
func ScrapBulldogJobs(path string) []ListingData {
	listings := []ListingData{}

	doc := CreateGoQuery(path)

	doc.Find(bulldogJobsTitleClass).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s - found position: %s\n", bulldogJobs, s.Text())

		listing := ListingData{
			Page:  bulldogJobs,
			Title: s.Text(),
		}

		listings = append(listings, listing)
	})

	return listings
}
