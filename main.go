package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jancewicz/FresherScout/scripts"
)

var dataMap = map[string]PageDetails{
	"noFluffJobs": {
		Url: "https://nofluffjobs.com/pl/Golang",
	},
	"protocolIT": {
		Url: "https://theprotocol.it/praca?kw=golang",
	},
}

func main() {
	var wg sync.WaitGroup
	htmlPathChan := make(chan string)

	fmt.Println("Lets GO scout!")

	for key, val := range dataMap {
		wg.Add(1)
		go func(key, url string) {
			defer wg.Done()
			htmlPath := ScrapPage(key, url)
			htmlPathChan <- htmlPath
		}(key, val.Url)
	}

	// channel is closed after all scraps are done
	go func() {
		wg.Wait()
		close(htmlPathChan)
	}()

	var processingWg sync.WaitGroup

	for htmlFilePath := range htmlPathChan {
		processingWg.Add(1)

		go func(path string) {
			defer processingWg.Done()

			switch {
			case strings.Contains(path, "noFluffJobs"):
				scripts.ScrapNFJHTML(path)
				fmt.Printf("NFJ scraping completed for: %s\n", path)
			case strings.Contains(path, "protocol"):
				scripts.ScrapProtocol(path)
				fmt.Printf("Protocol scraping completed for: %s\n", path)
			default:
				fmt.Printf("Unknown file : %s\n", path)
			}

			fmt.Printf("Completed processing file: %s\n", path)
		}(htmlFilePath)

		fmt.Printf("Started processing file: %s\n", htmlFilePath)
	}

	processingWg.Wait()
	fmt.Println("Scouting done!")
}
