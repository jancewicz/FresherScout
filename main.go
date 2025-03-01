package main

import (
	"fmt"

	"github.com/jancewicz/FresherScout/scripts"
)

func main() {
	fmt.Println("Lets GO scout!")
	htmlPathChan := make(chan string)
	scrapDoneChan := make(chan struct{})

	// reads data from map for each platform
	for key, val := range dataMap {
		go ScrapPage(key, val.Url, htmlPathChan)

		htmlFilePath := <-htmlPathChan

		go scripts.ScrapNFJHTML(htmlFilePath, scrapDoneChan)
		<-scrapDoneChan
		fmt.Printf("Scrapping file: %s completed\n", htmlFilePath)
	}

	fmt.Println("Scouting done!")
}
