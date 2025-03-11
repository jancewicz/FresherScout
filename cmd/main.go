package main

import (
	"fmt"
	"sync"

	"github.com/jancewicz/FresherScout/scripts"
)

func main() {
	var wg sync.WaitGroup
	pagaDataChan := make(chan PageData)

	fmt.Println("Lets GO scout!")

	for key, val := range scripts.DataMap {
		wg.Add(1)
		go func(key, url string) {
			defer wg.Done()
			htmlPath := ScrapPage(key, url)
			pagaDataChan <- htmlPath
		}(key, val.Url)
	}

	go func() {
		wg.Wait()
		close(pagaDataChan)
	}()

	var processingWg sync.WaitGroup

	for pageData := range pagaDataChan {
		processingWg.Add(1)

		selector := scripts.GetSelector(pageData.Name)

		go func(data PageData) {
			defer processingWg.Done()

			if err := scripts.Execute(scripts.ScrapHTMLFile, data.HTML, data.Name, selector, data.CSV); err != nil {
				fmt.Printf("Error during scrapping: %s, err: %v", data.Name, err)
			}

			fmt.Printf("Completed processing file: %s\n", data.HTML)
		}(pageData)

		fmt.Printf("Started processing file: %s\n", pageData.HTML)
	}

	processingWg.Wait()
	fmt.Println("Scouting done!")
}
