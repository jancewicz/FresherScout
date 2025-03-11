package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type PageData struct {
	HTML string
	CSV  string
	Name string
}

//	 Using scrapingBee api function encodes page addres and save its HTML to separate directory
//		name: name of scrapped page, needed for directories and files creation
//		addr: address of scrapped page
func ScrapPage(name, addr string) PageData {
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

	return PageData{
		HTML: fmt.Sprintf("files/%s/%s.html", name, name),
		CSV:  fmt.Sprintf("files/%s/%s.csv", name, name),
		Name: name,
	}
}
