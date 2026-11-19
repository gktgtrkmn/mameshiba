package main

import (
	"fmt"
	"os"
	"time"
)

func liveBubble(greeting, fact string, delay time.Duration) {
	lines, w := bubbleLines(greeting, fact)
	height := len(lines) + 2 + len(beanArt)
	fmt.Println(frame(reveal(lines, 0), w))
	for n := 1; n <= runeCount(lines); n++ {
		time.Sleep(delay)
		fmt.Printf("\x1b[%dA", height)
		fmt.Println(frame(reveal(lines, n), w))
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
		liveBubble(greeting, fact, typeDelay)
	default:
		fmt.Println(bubble(greeting, fact))
	}
	return nil
}
