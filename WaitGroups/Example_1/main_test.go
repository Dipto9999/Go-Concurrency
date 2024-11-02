package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printString(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	var wg sync.WaitGroup

	wg.Add(1)
	go printString("Test_1", &wg)
	wg.Wait()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "Test_1") {
		t.Errorf(
			"Expected Value: %s\n Actual Value: %s\n",
			"Test_1", output,
		)
	}
}
