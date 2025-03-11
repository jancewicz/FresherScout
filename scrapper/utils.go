package scrapper

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type PageDetails struct {
	Name         string
	Url          string
	CSSselectors CSSselectors
	Directory    string
}

type CSSselectors struct {
	JobTitle string
}

var DataMap = map[string]PageDetails{
	"noFluffJobs": {
		Url: "https://nofluffjobs.com/pl/Golang",
		CSSselectors: CSSselectors{
			JobTitle: ".posting-title__position.ng-star-inserted",
		},
	},
	"protocolIT": {
		Url: "https://theprotocol.it/praca?kw=golang",
		CSSselectors: CSSselectors{
			JobTitle: "#offer-title",
		},
	},
	"bulldogJobs": {
		Url: "https://bulldogjob.pl/companies/jobs/s/skills,Go",
		CSSselectors: CSSselectors{
			JobTitle: ".JobListItem_item__title__278xz",
		},
	},
}

var JobPoistions = []string{"Junior", "junior", "Trainee", "trainee", "Intern", "intern"}
var JuniorJobsFound = []ListingData{}

type ListingData struct {
	Page  string
	Title string
}

// Execute scrapper for given html file path
func Execute(scrapper func(htmlPath, pageName, jobTitleSelector string) []ListingData, htmlPath, pageName, jobTitleSelector, csvPath string) error {
	listings := scrapper(htmlPath, pageName, jobTitleSelector)

	if err := SaveListings(csvPath, listings); err != nil {
		return err
	}

	CheckPositions(csvPath)
	return nil
}

func CreateGoQuery(path string) *goquery.Document {
	htmlFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open the file: %v", err)
		return nil
	}
	defer htmlFile.Close()

	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		fmt.Printf("Cannot read HTML file: %v\n", err)
		return nil
	}

	return doc
}

func GetSelector(name string) string {
	_, exist := DataMap[name]

	if !exist {
		fmt.Println("Key does not exist")
	}
	return DataMap[name].CSSselectors.JobTitle
}
