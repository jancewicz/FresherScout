package main

import (
	"fmt"
	"sync"

	scrapp "github.com/jancewicz/FresherScout/testScripts"
)

func main() {
	fmt.Println("Lets GO scout!")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		scrapp.ScrapFirst()
	}()

	go func() {
		defer wg.Done()
		scrapp.ScrapSecond()
	}()

	wg.Wait()
	fmt.Println("Scouting done!")
}
