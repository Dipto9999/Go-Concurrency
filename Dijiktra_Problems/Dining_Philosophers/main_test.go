package main

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	for i := 0; i < 10; i++ {
		paid = []string{}
		main()

		if len(paid) != len(philosophers) {
			t.Errorf(
				"Expected Length of Paid Philosophers: %d"+
					"Actual Length of Paid Philosophers: %d",
				len(philosophers), len(paid),
			)
		}
	}
}

func Test_main_varydelays(t *testing.T) {
	var tests = map[string]time.Duration{
		"0 s":    0 * time.Second,
		"0.25 s": 250 * time.Millisecond,
		"0.5 s":  500 * time.Millisecond,
	}

	for _, std_delay := range tests {
		paid = []string{}
		THINK_TIME = 1 * std_delay
		EAT_TIME = 3 * std_delay
		SLEEP_TIME = 5 * std_delay

		main()
		if len(paid) != len(philosophers) {
			t.Errorf(
				"Test (): %s"+
					"Expected Length of Paid Philosophers: %d"+
					"Actual Length of Paid Philosophers: %d",
				std_delay, len(philosophers), len(paid),
			)
		}
	}

}
