package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

const WEEKS_IN_YEAR int = 52

type Income struct {
	Source string
	Amount int
}

var incomes []Income

func main() {
	var bankBalance int
	var balance sync.Mutex

	incomes = []Income{
		{Source: "Full Time Engineer @Tesla Inc.", Amount: 2000},
		{Source: "Part Time Math/Programming Tutor", Amount: 250},
		{Source: "TFSA Investments", Amount: 100},
		{Source: "Birthday Gifts", Amount: 5},
		{Source: "Savings Account", Amount: 2},
	} // Define Weekly Revenue

	// Bank Balance
	fmt.Printf("Initial Account Balance: $%d.00\n", bankBalance) // Starting Values

	wg.Add(len(incomes))
	// Determine How Much Made in 1 Year; Keep Running Total
	for i, income := range incomes {
		go func(i int, income Income) {
			defer wg.Done()

			for week := 1; week <= WEEKS_IN_YEAR; week++ {
				balance.Lock()
				bankBalance += income.Amount
				balance.Unlock()
				// fmt.Printf("Week %d: Earned $%d.00 from %s\n", week, income.Amount, income.Source)
			}

		}(i, income)
	}
	wg.Wait()

	fmt.Printf("Final Bank Balance: $%d.00", bankBalance) // Final Balance
}
