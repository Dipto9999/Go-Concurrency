/*
 **********************************************************************************************
 * Simple Demonstration of Solving Sleeping Barber dilemma, a classic Computer Science Problem.
 **********************************************************************************************
 * This illustrates the complexities that arise when there are multiple operating system processes.
 * Here, we have a finite number of barbers, a finite number of seats in a waiting room, a fixed length of
 * time the barbershop is open, and clients arriving at (roughly) regular intervals. When a barber has
 * nothing to do, he/she checks the waiting room, and if 1+ clients are there, a haircut takes place. Otherwise,
 * the barber goes to sleep until a new client arrives.
 **********************************************************************************************
 * The rules are as follows:
 * 	- when barber finishes a haircut, he inspects the waiting room to see if any waiting customers
 * 		- if 0 customers, the barber falls asleep in chair
 *  	- a customer must wake the barber if he is asleep
 * 	- barber cannot leave until waiting room is empty
 * 	- shop can stop accepting new clients at closing time
 * 		- after shop is closed and no clients left in waiting area, barber goes home
 **********************************************************************************************
 * Can solve without the use of semaphores (mutexes)
 */

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed Random Number Generator

	color.Black("Sleeping Barber Problem")
	color.Black("-----------------------")

	// Create Channels
	clientsChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	barbershop := Barbershop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientsChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	} // Open Barbershop

	color.Cyan("Shop is Open for the Day!")

	// Add Barbers
	barbershop.addBarber("Frank")
	barbershop.addBarber("John")
	barbershop.addBarber("Jake")
	barbershop.addBarber("Maria")
	barbershop.addBarber("Jane")

	// Start Barbershop as Go Routine
	doorClosed := make(chan bool)
	shopClosed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		doorClosed <- true
		barbershop.closeShopForDay()
		shopClosed <- true
	}()

	// Add Clients
	i := 1
	go func() {
		for {
			arrivalTime := rand.Int() % (2 * arrivalRate) // Get Random Number with Avg Arrival Rate
			select {
			case <-doorClosed:
				return
			case <-time.After(time.Duration(arrivalTime) * time.Millisecond):
				barbershop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	// Block Until Barbershop is Closed
	<-shopClosed
}
