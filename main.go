package main

import (
	"fmt"
	"sync"

	"github.com/jancewicz/FresherScout/scripts"
)

var dataMap = map[string]PageDetails{
	"noFluffJobs": {
		Url:         "https://nofluffjobs.com/pl/Golang",
		CSSselector: ".posting-title__position.ng-star-inserted",
	},
	"protocolIT": {
		Url:         "https://theprotocol.it/praca?kw=golang",
		CSSselector: "#offer-title",
	},
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	htmlPathChan := make(chan string)
	scrapDoneChan := make(chan struct{})

	fmt.Println("Lets GO scout!")

	for key, val := range dataMap {
		wg.Add(1)
		go func(key, url string) {
			defer wg.Done()
			ScrapPage(key, url, htmlPathChan)
		}(key, val.Url)
	}

	// channel is closed after all scraps are done
	go func() {
		wg.Wait()
		close(htmlPathChan)
	}()

	for htmlFilePath := range htmlPathChan {
		wg.Add(2)

		go func(path string) {
			defer wg.Done()
			scripts.ScrapNFJHTML(path, scrapDoneChan)
		}(htmlFilePath)

		go func(path string) {
			defer wg.Done()
			scripts.ScrapProtocol(path, scrapDoneChan)
		}(htmlFilePath)

		<-scrapDoneChan
		fmt.Printf("Scrapping file: %s completed\n", htmlFilePath)
	}

	wg.Wait()
	fmt.Println("Scouting done!")
}
