package main

import (
	"fmt"
	"os"
	"time"
)

func liveBubble(r liveRenderer, t typing, delay time.Duration) {
	fmt.Print(r.start(t))
	for !t.done() {
		time.Sleep(delay)
		t = t.next()
		fmt.Print(r.step(t))
	}
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	return (err == nil) && (fi.Mode()&os.ModeCharDevice != 0)
}

func runCLI(fetch func() (string, error), live, plain bool) error {
	fact, err := fetch()
	if err != nil {
		return fmt.Errorf("could not fetch a fact: %w", err)
	}
	switch {
	case plain || !isTerminal(os.Stdout):
		fmt.Println(fact)
	case live:
		liveBubble(cliLive{}, newTyping(greeting, fact), typeDelay)
	default:
		fmt.Println(bubble(greeting, fact))
	}
	return nil
}
