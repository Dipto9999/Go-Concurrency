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

var KITCHEN_OUTCOMES = map[string]int{
	"Fire":                1,
	"Pizza_Burned":        2,
	"Pizza_NoIngredients": 4,
	"Pizza_Success":       5,
}

type Pizzeria struct {
	pizzaOrders chan PizzaOrder
	quit        chan struct{}
}

type PizzaOrder struct {
	id      int
	message string
	success bool
}

func (p *Pizzeria) Close() {
	close(p.quit)
}

/*
 * Attempt to Make a Pizza:
 * 	1) Generate a Random Number Between 1-10
 *  	i) If outcome <= 2 : Pizza was Burned
 *  	ii) If outcome <= 4 : Pizza Missing Ingredients
 *  	iii) If outcome >= 5 : Pizza Successfully Made
 *	2) Pizzas Will Take Different Amounts of Time to Make.
 */
func makePizza(id int, quit chan struct{}) *PizzaOrder {
	id++
	if id <= NumberOfPizzas {
		fmt.Printf("Received Pizza Order ID#%d!\n", id)

		delay := rand.Intn(5) + 1
		fmt.Printf("Making Pizza ID#%d. It will take %d s...\n", id, delay)
		time.Sleep(time.Duration(delay) * time.Second) // Delay to Simulate Cooking

		select {
		case <-quit:
			return nil // Stop Making Pizza.
		default:
			msg, outcome := "", rand.Intn(9)+1
			if outcome <= KITCHEN_OUTCOMES["Pizza_Burned"] {
				msg = fmt.Sprintf("*** Pizza ID#%d Burned in the Oven!", id)
			} else if outcome <= KITCHEN_OUTCOMES["Pizza_NoIngredients"] {
				msg = fmt.Sprintf("*** Ran Out of Ingredients for Pizza ID#%d!", id)
			} else {
				msg = fmt.Sprintf("Pizza Order ID#%d is Ready!", id)
			}

			success := false
			if outcome < KITCHEN_OUTCOMES["Pizza_Success"] {
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
	}
	return nil
}

/*
 * Keep Track of Which Pizza Being Made.
 * Run Forever or Until We Receive a Quit Notification.
 */
func (pizzeria *Pizzeria) kitchen(wg *sync.WaitGroup) {
	defer wg.Done()

	var i = 0
	for {
		// Check for Quit Signal.
		select {
		case <-pizzeria.quit:
			// Close Channels
			close(pizzeria.pizzaOrders)
			return
		default:
			currentPizza := makePizza(i, pizzeria.quit)
			if currentPizza != nil {
				i = currentPizza.id
				pizzeria.pizzaOrders <- *currentPizza // Send Pizza Order to Customer.
			} else {
				// All Pizzas Have Been Attempted.
				close(pizzeria.pizzaOrders) // Close pizzaOrders Channel to signal completion.
				return
			}
		}
	}
}

/*
 * Tourists Expecting to Consume NumberOfPizzas.
 */
func (pizzeria *Pizzeria) tourists(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range pizzeria.pizzaOrders {
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
func (pizzeria *Pizzeria) maintenance(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second) // Check Every 1-4 s.

		select {
		case <-pizzeria.quit: // Exit if Quit Signal Received.
			return
		default:
			outcome := rand.Intn(9) + 1
			if outcome <= KITCHEN_OUTCOMES["Fire"] {
				// Simulate Early Shutdown.
				color.HiRed("*** Panicking! Fire in Pizzeria!***\n***Everybody Leave the Premises! Shutting Down for the Day ***")
				pizzeria.Close()

				return
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed Random Number Generator

	pizzeria := &Pizzeria{
		pizzaOrders: make(chan PizzaOrder),
		quit:        make(chan struct{}),
	} // Create PizzaJobs

	color.Cyan("The Pizzeria is Open for Business!")
	color.Cyan("----------------------------------")

	wg.Add(1)
	go pizzeria.kitchen(&wg)
	wg.Add(1)
	go pizzeria.tourists(&wg)
	wg.Add(1)
	go pizzeria.maintenance(&wg)

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
