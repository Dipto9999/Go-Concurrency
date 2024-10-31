package main

import (
	"fmt"
	"sync"
)

func printString(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(s)
}

func main() {
	var wg sync.WaitGroup

	weekdays := []string{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
	}

	wg.Add(len(weekdays))
	fmt.Println("Wait Group: ", wg)

	for i, day := range weekdays {
		go printString(fmt.Sprintf("%d: %s", i, day), &wg)
	}
	wg.Wait()

	for _, phrase := range []string{"Muntakim", "Hello World!", "Mahir"} {
		wg.Add(1)
		go printString(phrase, &wg)
	}
	wg.Wait()

	fmt.Println("Wait Group: ", wg)
}
