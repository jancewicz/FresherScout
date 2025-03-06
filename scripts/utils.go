package scripts

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var JobPoistions = []string{"Junior", "junior", "Trainee", "trainee", "Intern", "intern"}
var JuniorJobsFound = []ListingData{}

type ListingData struct {
	Page  string
	Title string
}

var csvHeaders = []string{
	"Position",
}

func Execute(scrapper func(htmlPath string) []ListingData, htmlPath, csvPath string) error {
	listings := scrapper(htmlPath)

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

func SaveListings(path string, listings []ListingData) error {
	csvFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Cannot create file", err)
		return err
	}
	defer csvFile.Close()

	if err := WriteContent(listings, csvFile); err != nil {
		fmt.Println("error during appending data to csv file: ", err)
		return err
	}

	return nil
}

func WriteContent(listings []ListingData, csvFile *os.File) error {
	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	err := csvwriter.Write(csvHeaders)
	if err != nil {
		return err
	}

	for _, offer := range listings {
		row := []string{offer.Title}
		if err := csvwriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// WIP

func CheckPositions(csvPath string) bool {
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("file does not exist")
		return false
	}

	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Cannot open file", err)
		return false
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Cannot read file", err)
		return false
	}

	for _, row := range data {
		flatRow := strings.Join(row, " ")
		if ContainAny(flatRow, JobPoistions) {
			fmt.Printf("Job position found: %s\n", flatRow)
			return true
		}
	}

	fmt.Println("Cannot find junior position")
	return false
}

func ContainAny(line string, positions []string) bool {
	for _, position := range positions {
		if strings.Contains(line, position) {
			return true
		}
	}
	return false
}
