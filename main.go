package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const typeDelay = 30 * time.Millisecond

func main() {
	live := flag.Bool("live", false, "print fact as stream")
	plain := flag.Bool("plain", false, "print only the fact, no bubble")
	tui := flag.Bool("tui", false, "fullscreen mode, mameshiba says a fact periodically")
	interval := flag.Duration("interval", 15*time.Second, "time between facts in tui mode")
	flag.Parse()

	run := func() error { return runCLI(fetchFact, *live, *plain) }
	if *tui {
		run = func() error { return runTUI(fetchFact, *interval) }
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
