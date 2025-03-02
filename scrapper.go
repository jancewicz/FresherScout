package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var jobPoistions = []string{"Junior", "junior", "Trainee", "trainee", "Intern", "intern"}

type PageDetails struct {
	Url         string
	CSSselector string
	Directory   string
}

//	 Using scrapingBee api function encodes page addres and save its HTML to separate directory
//		name: name of scrapped page, needed for directories and files creation
//		addr: address of scrapped page
func ScrapPage(name, addr string, c chan string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error occured during .env file loading")
		return
	}
	apiKey := os.Getenv("API_KEY")
	encodedUrl := url.QueryEscape(addr)

	client := http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://app.scrapingbee.com/api/v1/?api_key=%s&url=%s", apiKey, encodedUrl), nil)
	if err != nil {
		fmt.Printf("error occured during request creation: %v", err)
	}

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Printf("error occured during form parsing: %v", parseFormErr)
		return
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("error occured during request: %v", err)
		return
	}

	responseBody, _ := io.ReadAll(response.Body)
	dirPath := fmt.Sprintf("files/%s", name)
	os.MkdirAll(dirPath, os.ModePerm)

	file, err := os.Create(fmt.Sprintf("files/%s/%s.html", name, name))
	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}
	defer file.Close()
	file.Write(responseBody)

	defer response.Body.Close()

	c <- fmt.Sprintf("files/%s/%s.html", name, name)
}

func CheckPositions(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("file does not exist")
		return false
	}

	file, err := os.Open(path)
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
		if containAny(flatRow, jobPoistions) {
			fmt.Println("Job position found!")
			return true
		}
	}

	fmt.Println("Cannot find junior position")
	return false
}

func containAny(line string, positions []string) bool {
	for _, position := range positions {
		if strings.Contains(line, position) {
			return true
		}
	}
	return false
}
