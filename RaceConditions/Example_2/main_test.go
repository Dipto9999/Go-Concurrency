package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	weekly_income := 0
	for _, curr_income := range incomes {
		weekly_income += curr_income.Amount
	}
	annual_income := WEEKS_IN_YEAR * weekly_income

	if !strings.Contains(output, fmt.Sprintf("$%d.00", annual_income)) {
		t.Errorf(
			"Expected Annual Income: $%d.00\n Actual Output: %s\n",
			annual_income, output,
		)
	}
}
