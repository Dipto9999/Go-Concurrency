package main

import (
	"time"

	"github.com/fatih/color"
)

type Barbershop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (barbershop *Barbershop) addBarber(barber string) {
	barbershop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Cyan("%s Checking for Clients in Waiting Room...", barber)

		for {
			// If No Clients, Barber Goes to Sleep
			if len(barbershop.ClientsChan) == 0 {
				color.Black("Nothing to Do...%s Takes a Nap...Zzz", barber)
				isSleeping = true
			}

			client, shopIsOpen := <-barbershop.ClientsChan

			if shopIsOpen {
				if isSleeping {
					color.Green("%s -> Wakes up %s", client, barber)
					isSleeping = false
				}
				barbershop.cutHair(barber, client)
			} else {
				barbershop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (barbershop *Barbershop) cutHair(barber string, client string) {
	color.Green("%s Giving %s a Haircut...", barber, client)
	time.Sleep(barbershop.HairCutDuration)
	color.Green("%s Finished Giving %s a Haircut!", barber, client)
}

func (barbershop *Barbershop) sendBarberHome(barber string) {
	color.HiBlack("*%s Going Home*", barber)
	barbershop.BarbersDoneChan <- true
}

func (barbershop *Barbershop) closeShopForDay() {
	color.HiBlack("Closing Shop for the Day")

	close(barbershop.ClientsChan)
	barbershop.Open = false

	for a := 1; a <= barbershop.NumberOfBarbers; a++ {
		<-barbershop.BarbersDoneChan
	}

	close(barbershop.BarbersDoneChan)
	color.HiBlack("---------------------------------------------------------")
	color.HiBlack("Barbershop is Closed for the Day. Everyone has Gone Home.")
}

func (barbershop *Barbershop) addClient(client string) {
	color.Blue("%s Arrives!", client)

	if barbershop.Open {
		select {
		case barbershop.ClientsChan <- client:
			color.Blue("%s Takes a Seat in Waiting Room", client)
		default:
			color.Red("Waiting Room is Full! %s Leaves...", client)
		}
	} else {
		color.Red("Shop is Already Closed! %s Leaves...", client)
	}
}
