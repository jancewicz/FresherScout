package scripts

import (
	"encoding/csv"
	"os"
)

type ListingData struct {
	Title string
}

var csvHeaders = []string{
	"Position",
}

func WriteListingsToCSV(listings []ListingData, csvFile *os.File) error {
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
