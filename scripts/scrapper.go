package scripts

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func ScrapHTMLFile(path, pageName, jobTitle string) []ListingData {
	var listings = []ListingData{}

	doc := CreateGoQuery(path)

	doc.Find(jobTitle).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s - found position: %s\n", pageName, s.Text())

		listing := ListingData{
			Page:  pageName,
			Title: s.Text(),
		}

		listings = append(listings, listing)
	})

	return listings
}
