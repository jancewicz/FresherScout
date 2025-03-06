package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type PageDetails struct {
	Url         string
	CSSselector string
	Directory   string
}

type FilePaths struct {
	HTML string
	CSV  string
}

func GenerateCSVPath(name string) string {
	return fmt.Sprintf("files/%s/%s.csv", name, name)
}

//	 Using scrapingBee api function encodes page addres and save its HTML to separate directory
//		name: name of scrapped page, needed for directories and files creation
//		addr: address of scrapped page
func ScrapPage(name, addr string) FilePaths {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error occured during .env file loading")
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
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("error occured during request: %v", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	dirPath := fmt.Sprintf("files/%s", name)
	os.MkdirAll(dirPath, os.ModePerm)

	file, err := os.Create(fmt.Sprintf("files/%s/%s.html", name, name))
	if err != nil {
		fmt.Println("Cannot create file", err)
	}
	defer file.Close()
	file.Write(responseBody)

	defer response.Body.Close()

	return FilePaths{
		HTML: fmt.Sprintf("files/%s/%s.html", name, name),
		CSV:  fmt.Sprintf("files/%s/%s.csv", name, name),
	}
}

func containAny(line string, positions []string) bool {
	for _, position := range positions {
		if strings.Contains(line, position) {
			return true
		}
	}
	return false
}
