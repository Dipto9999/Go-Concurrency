/*
 * The Dining Philosophers problem is well known in Computer Science circles.
 * 5 philosophers, numbered from 0-4, live in a house where the table is laid for them;
 * Each philosopher has their own place at the table. Their only difficulty - besides those of philosphy - is that the dish
 * served is a very difficult kind of spaghetti which has to be eaten with 2 forks. There are 2 forks next to each plate, so that
 * presents no difficulty. As a consequence however, this means that no 2 neighbours may be eating simulatenously,
 * since there are 5 philosophers and 5 forks.
 *
 * This is a simple implementation of Dijkstra's solution in the "Dining Philosophers" dilemma.
 */
package main

import (
	"sync"
	"time"

	"github.com/fatih/color"
)

type Philosopher struct {
	name      string
	leftFork  int
	rightFork int
}

var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

const HUNGER = 100
const THINK_TIME = 300 * time.Millisecond
const EAT_TIME = 100 * time.Millisecond
const SLEEP_TIME = 100 * time.Millisecond

var seated = &sync.WaitGroup{}
var hungry = &sync.WaitGroup{}
var forks = make(map[int]*sync.Mutex)

var DINING_OPTIONS = map[string]int{
	"UNINTERRUPTED": 1,
	"SYMMETRICAL":   2,
}

func eat(philosopher Philosopher, forks map[int]*sync.Mutex) {
	color.Green("Philosopher %s is Eating...", philosopher.name)
	time.Sleep(EAT_TIME)
	forks[philosopher.leftFork].Unlock()
	forks[philosopher.rightFork].Unlock()
	color.Green("Philosopher %s Has Released Forks %d and %d", philosopher.name, philosopher.leftFork, philosopher.rightFork)
}

func think(philosopher Philosopher) {
	color.Black("Philosopher %s is Thinking...", philosopher.name)
	time.Sleep(THINK_TIME)
}

func dine_uninterrupted(philosopher Philosopher, forks map[int]*sync.Mutex) {
	acquired := false
	for i := HUNGER; i > 0; i-- {
		acquired = false

		// Try to Acquire Both Necessary Forks.
		forks[philosopher.leftFork].Lock()
		color.Green("Philosopher %s Takes Fork %d on Left.", philosopher.name, philosopher.leftFork)
		if forks[philosopher.rightFork].TryLock() { // Non-Blocking
			acquired = true
			color.Green("Philosopher %s Takes Fork %d on Right.", philosopher.name, philosopher.rightFork)
		} else {
			forks[philosopher.leftFork].Unlock()
			color.Red("Philosopher %s Releases Fork %d on Left.", philosopher.name, philosopher.leftFork)
		}

		if acquired {
			eat(philosopher, forks)
		}
		think(philosopher)
	}

	color.HiGreen("Philosopher %s Has Finished Dining!", philosopher.name)
}

func dine_symmetrical(philosopher Philosopher, forks map[int]*sync.Mutex) {
	for i := HUNGER; i > 0; i-- {
		if philosopher.leftFork%2 == 0 {
			forks[philosopher.leftFork].Lock()
			color.Green("Philosopher %s Takes Fork %d on Left.", philosopher.name, philosopher.leftFork)
			forks[philosopher.rightFork].Lock()
			color.Green("Philosopher %s Takes Fork %d on Right.", philosopher.name, philosopher.rightFork)
		} else {
			forks[philosopher.rightFork].Lock()
			color.Green("Philosopher %s Takes Fork %d on Right.", philosopher.name, philosopher.rightFork)
			forks[philosopher.leftFork].Lock()
			color.Green("Philosopher %s Takes Fork %d on Left.", philosopher.name, philosopher.leftFork)
		}

		eat(philosopher, forks)
		think(philosopher)
	}

	color.HiGreen("Philosopher %s Has Finished Dining!", philosopher.name)
}

func dine(philosopher Philosopher, forks map[int]*sync.Mutex, how int) {
	defer hungry.Done()

	// Seat Philosopher at Table.
	color.Black("Philosopher %s Seated at Table.", philosopher.name)
	seated.Done()

	seated.Wait()

	if how == DINING_OPTIONS["SYMMETRICAL"] {
		dine_symmetrical(philosopher, forks)
	} else {
		dine_uninterrupted(philosopher, forks)
	}
}

func main() {
	seated.Add(len(philosophers))
	hungry.Add(len(philosophers))
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	color.Cyan("Dining Philosopher's Problem")
	color.Cyan("----------------------------")
	color.Cyan("Table is Empty")

	// Start Meal
	for i := 0; i < len(philosophers); i++ {
		go dine(philosophers[i], forks, DINING_OPTIONS["SYMMETRICAL"]) // Fire off Go Routine for Current Philosopher
	}
	hungry.Wait()

	color.Cyan("----------------------------")
	color.Cyan("Table is Empty")
}
