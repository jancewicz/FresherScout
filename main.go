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
	htmlPathChan := make(chan FilePaths)

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

		go func(paths FilePaths) {
			defer processingWg.Done()

			switch {
			case strings.Contains(paths.HTML, "noFluffJobs"):
				scripts.ScrapNFJHTML(paths.HTML)
				fmt.Printf("NFJ scraping completed for: %s\n", paths.HTML)
			case strings.Contains(paths.HTML, "protocol"):
				scripts.ScrapProtocol(paths.HTML)
				fmt.Printf("Protocol scraping completed for: %s\n", paths.HTML)
			default:
				fmt.Printf("Unknown file : %s\n", paths.HTML)
			}

			fmt.Printf("Completed processing file: %s\n", paths.HTML)

			CheckPositions(paths.CSV)
		}(htmlFilePath)

		fmt.Printf("Started processing file: %s\n", htmlFilePath)
	}

	processingWg.Wait()
	fmt.Println("Scouting done!")

}
