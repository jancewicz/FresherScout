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
	"bulldogJobs": {
		Url: "https://bulldogjob.pl/companies/jobs/s/skills,Go",
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
				if err := scripts.Execute(scripts.ScrapNFJHTML, paths.HTML, paths.CSV); err != nil {
					fmt.Println("Error during scrapping noFluffJobs")
				}
				fmt.Printf("NFJ scraping completed for: %s\n", paths.HTML)
			case strings.Contains(paths.HTML, "protocol"):
				if err := scripts.Execute(scripts.ScrapProtocol, paths.HTML, paths.CSV); err != nil {
					fmt.Println("Error during scrapping noFluffJobs")
				}
				fmt.Printf("Protocol scraping completed for: %s\n", paths.HTML)
			case strings.Contains(paths.HTML, "bulldog"):
				if err := scripts.Execute(scripts.ScrapBulldogJobs, paths.HTML, paths.CSV); err != nil {
					fmt.Println("Error during scrapping noFluffJobs")
				}
				fmt.Printf("BulldogJobs scraping completed for: %s\n", paths.HTML)
			default:
				fmt.Printf("Unknown file : %s\n", paths.HTML)
			}

			fmt.Printf("Completed processing file: %s\n", paths.HTML)
		}(htmlFilePath)

		fmt.Printf("Started processing file: %s\n", htmlFilePath)
	}

	processingWg.Wait()
	fmt.Println("Scouting done!")
}
