package scripts

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const protocolIT = "protocolIT"

var titleIDName = "#offer-title"

// Path: path to  html file
func ScrapProtocol(path string) []ListingData {
	listings := []ListingData{}

	doc := CreateGoQuery(path)

	doc.Find(titleIDName).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s - found position: %s\n", protocolIT, s.Text())
		listings = append(listings, ListingData{
			Title: s.Text(),
		})
	})

	return listings
}
