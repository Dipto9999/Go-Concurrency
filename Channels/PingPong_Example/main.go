package main

import (
	"fmt"
	"strings"
)

/*
 * Takes String from Ping Channel, Converts to Uppercase and Appends
 * Exclamation Mark. Converted String is Sent to Pong Channel.
 */
func shout(ping <-chan string, pong chan<- string) {
	for {
		s := <-ping // Read from Pong Channel (Blocking)
		pong <- fmt.Sprintf("%s!", strings.ToUpper(s))
	}
}

func main() {
	// Create Channels
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type Something and Press Enter / Q to Quit.")
	for {
		fmt.Printf("-> ")
		var userInput string
		_, _ = fmt.Scanln(&userInput) // Scan User Input

		if strings.ToUpper(userInput) == "Q" {
			break // Quit Program
		} else {
			ping <- userInput  // Send to Ping Channel
			response := <-pong // Wait for Response from Pong Channel
			fmt.Println("Response : ", response)
		}
	}

	fmt.Println("All Done. Closing Channels")
	// Close Channels
	close(ping)
	close(pong)
}
