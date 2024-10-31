package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	wg.Add(1)
	go updateMessage("Test_UpdateMessage", &wg)
	wg.Wait()

	if msg != "Test_UpdateMessage" {
		t.Errorf(
			"Expected Value: %s\nActual Value: %s\n",
			"Test_UpdateMessage", msg,
		)
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	msg = "Test_printMessage"
	printMessage()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "Test_printMessage") {
		t.Errorf(
			"Expected Value: %s\nActual Value: %s\n",
			"Test_UpdateMessage", output,
		)
	}
}

func test_main(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	for _, curr_string := range []string{
		"Hello, universe!", "Hello, cosmos!", "Hello, world!",
	} {
		if !strings.Contains(output, curr_string) {
			t.Errorf(
				"Expected Value: %s\nActual Value: %s\n",
				"Test_UpdateMessage", output,
			)
		}
	}
}
