package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var wg sync.WaitGroup
var pizzasMade, pizzasFailed int

var PIZZA_OUTCOMES = map[string]int{
	"Burned":              2,
	"Missing_Ingredients": 4,
	"Success":             5,
}

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	id      int
	message string
	success bool
}

func (p *Producer) Close() error {
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

		outcome := rand.Intn(9) + 1
		if outcome < PIZZA_OUTCOMES["Success"] {
			pizzasFailed++
		} else {
			pizzasMade++
		}

		delay := rand.Intn(5) + 1
		fmt.Printf("Making Pizza ID#%d. It will take %d s...\n", id, delay)
		time.Sleep(time.Duration(delay) * time.Second) // Delay to Simulate Cooking

		msg := ""
		success := false
		if outcome <= PIZZA_OUTCOMES["Burned"] {
			msg = fmt.Sprintf("*** Pizza ID#%d Burned in the Oven!", id)
		} else if outcome <= PIZZA_OUTCOMES["Missing_Ingredients"] {
			msg = fmt.Sprintf("*** Ran Out of Ingredients for Pizza ID#%d!", id)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza Order ID#%d is Ready!", id)
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
func pizzeria(pizzaMaker *Producer, wg *sync.WaitGroup) {
	defer wg.Done()

	var i = 0
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.id
			select {
			case pizzaMaker.data <- *currentPizza: // Send Pizza Order to Customer
			case quitChan := <-pizzaMaker.quit: // Listen for Quit Signal and Send Chan Error
				// Close Channels
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		} else {
			// All Pizzas Have Been Attempted.
			close(pizzaMaker.data) // Close Data Channel.
			return
		}
	}
}

func tourists(pizzaJob *Producer, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range pizzaJob.data {
		if i.success {
			color.Green(i.message)
			color.Green("Pizza Order ID#%d Out for Delivery!", i.id)
		} else {
			color.Red("Customer of Pizza Order ID#%d is Really Mad!", i.id)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed Random Number Generator

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	} // Create Producer

	color.Cyan("The Pizzeria is Open for Business!")
	color.Cyan("----------------------------------")

	wg.Add(1)
	go pizzeria(pizzaJob, &wg)
	wg.Add(1)
	go tourists(pizzaJob, &wg)
	wg.Wait()

	color.Cyan("Done Making Pizzas...")
	color.Cyan("----------------")
	color.Green("Successfully Made %d Pizzas", pizzasMade)
	color.Red("Failed to Make %d Pizzas", pizzasFailed)

	switch {
	case pizzasMade >= 9:
		color.Green("Amazing Day!")
	case pizzasMade >= 6:
		color.Green("Mostly Productive Day!")
	case pizzasMade == 5:
		color.Yellow("An Okay Day.")
	case pizzasMade >= 3:
		color.Red("Not a Good Day!")
	default:
		color.Red("Awful Day!")
	}
}
