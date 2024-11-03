package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 100

var wg sync.WaitGroup
var pizzasMade, pizzasFailed int

var PIZZA_OUTCOMES = map[string]int{
	"Kitchen_Fire":        1,
	"Burned":              2,
	"Missing_Ingredients": 4,
	"Success":             5,
}

type PizzaJobs struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	id      int
	message string
	success bool
}

func (p *PizzaJobs) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

/*
 * Attempt to Make a Pizza:
 * 	1) Generate a Random Number Between 1-10
 *  	i) If outcome <= 2 : Pizza was Burned
 *  	ii) If outcome <= 4 : Pizza Missing Ingredients
 *  	iii) If outcome >= 5 : Pizza Successfully Made
 *	2) Pizzas Will Take Different Amounts of Time to Make.
 */
func makePizza(id int) *PizzaOrder {
	id++
	if id <= NumberOfPizzas {
		fmt.Printf("Received Pizza Order ID#%d!\n", id)

		delay := rand.Intn(5) + 1
		fmt.Printf("Making Pizza ID#%d. It will take %d s...\n", id, delay)
		time.Sleep(time.Duration(delay) * time.Second) // Delay to Simulate Cooking

		msg, outcome := "", rand.Intn(9)+1
		if outcome <= PIZZA_OUTCOMES["Burned"] {
			msg = fmt.Sprintf("*** Pizza ID#%d Burned in the Oven!", id)
		} else if outcome <= PIZZA_OUTCOMES["Missing_Ingredients"] {
			msg = fmt.Sprintf("*** Ran Out of Ingredients for Pizza ID#%d!", id)
		} else {
			msg = fmt.Sprintf("Pizza Order ID#%d is Ready!", id)
		}

		success := false
		if outcome < PIZZA_OUTCOMES["Success"] {
			pizzasFailed++
		} else {
			success = true
			pizzasMade++
		}

		return &PizzaOrder{
			id:      id,
			message: msg,
			success: success,
		}
	}
	return nil
}

/*
 * Keep Track of Which Pizza Being Made.
 * Run Forever or Until We Receive a Quit Notification.
 */
func pizzeria(pizzaJobs *PizzaJobs, wg *sync.WaitGroup) {
	defer wg.Done()

	var i = 0
	for {
		// Check for Quit Signal.
		select {
		case quitChan := <-pizzaJobs.quit:
			// Close Channels
			close(pizzaJobs.data)
			close(quitChan)
			return
		default:
			currentPizza := makePizza(i)
			if currentPizza != nil {
				i = currentPizza.id
				pizzaJobs.data <- *currentPizza // Send Pizza Order to Customer.
			} else {
				// All Pizzas Have Been Attempted.
				close(pizzaJobs.data) // Close Data Channel to signal completion.
				return
			}
		}
	}
}

func tourists(pizzaJobs *PizzaJobs, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range pizzaJobs.data {
		if i.success {
			color.Green(i.message)
			color.Green("Pizza Order ID#%d Out for Delivery!", i.id)
		} else {
			color.Red("Customer of Pizza Order ID#%d is Really Mad!", i.id)
		}
	}
}

/*
 * Every Few Seconds, There is Risk of Kitchen Fire.
 * In this Event, Pizzeria Will Shut Down for the Day.
 */
func kitchen_activities(pizzaJobs *PizzaJobs, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second) // Check Every 1-4 s.

		select {
		case <-pizzaJobs.quit: // Exit if Quit Signal Received.
			return
		default:
			outcome := rand.Intn(9) + 1
			if outcome <= PIZZA_OUTCOMES["Kitchen_Fire"] {
				// Simulate Early Shutdown.
				color.HiRed("*** Panicking! Fire in Pizzeria!***\n***Everybody Leave the Premises! Shutting Down for the Day ***")
				pizzaJobs.Close()

				return
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed Random Number Generator

	pizzaJobs := &PizzaJobs{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	} // Create PizzaJobs

	color.Cyan("The Pizzeria is Open for Business!")
	color.Cyan("----------------------------------")

	wg.Add(1)
	go pizzeria(pizzaJobs, &wg)
	wg.Add(1)
	go tourists(pizzaJobs, &wg)
	wg.Add(1)
	go kitchen_activities(pizzaJobs, &wg)

	wg.Wait()

	color.Cyan("Done Making Pizzas...")
	color.Cyan("----------------")
	color.Green("Successfully Made %d Pizzas", pizzasMade)
	color.Red("Failed to Make %d Pizzas", pizzasFailed)

	switch {
	case pizzasMade >= 0.9*NumberOfPizzas:
		color.Green("Amazing Day!")
	case pizzasMade >= 0.6*NumberOfPizzas:
		color.Green("Mostly Productive Day!")
	case pizzasMade >= 0.5*NumberOfPizzas:
		color.Yellow("An Okay Day.")
	case pizzasMade >= 0.3*NumberOfPizzas:
		color.Red("Not a Good Day!")
	default:
		color.Red("Awful Day!")
	}
}
