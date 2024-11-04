package main

import (
	"fmt"
	"time"
)

func listenToChannel(ch chan int) {
	for {
		i := <-ch
		fmt.Println("Received ", i, " From Channel")
		time.Sleep(1 * time.Second) // Simulate Other Work Being Done
	}
}

func server_a(ch chan string) {
	for {
		time.Sleep(6 * time.Second)
		ch <- "From Server A"
	}
}

func server_b(ch chan string) {
	for {
		time.Sleep(3 * time.Second)
		ch <- "From Server B"
	}
}

func main() {
	fmt.Println("Select with Channels")
	fmt.Println("--------------------")

	channel_a := make(chan string)
	channel_b := make(chan string)
	channel_c := make(chan int, 10) // Buffered Channel

	go server_a(channel_a)
	go server_b(channel_b)
	go listenToChannel(channel_c)

	time.Sleep(5 * time.Second)

	for i := 0; i <= 100; i++ {
		fmt.Println("Sending ", i, " To Channel...")
		channel_c <- i
		fmt.Println("Sent ", i, " To Channel!")

		select {
		case s1 := <-channel_a:
			fmt.Println("Case 1:", s1)
		case s2 := <-channel_a:
			fmt.Println("Case 2:", s2)
		case s3 := <-channel_b:
			fmt.Println("Case 3:", s3)
		case s4 := <-channel_b:
			fmt.Println("Case 4:", s4)
		default:
			// Avoid Deadlock.
		}
	}
	fmt.Println("Done!")
	close(channel_a)
	close(channel_b)
	close(channel_c)

}
