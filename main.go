package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

const (
	factURL   = "https://uselessfacts.jsph.pl/api/v2/facts/random?language=en"
	greeting  = "ねえ、知ってる？\n"
	typeDelay = 30 * time.Millisecond
)

var beanArt = []string{
	"   ╲",
	"  (◕ ᴥ ◕ ᗜ   )",
}

func fetchFact() (string, error) {
	var body struct {
		Text string `json:"text"`
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(factURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Text == "" {
		return "", fmt.Errorf("no text field in response")
	}
	return body.Text, nil
}

func wrap(text string, width int) []string {
	var lines []string
	line := ""
	for w := range strings.FieldsSeq(text) {
		switch {
		case line == "":
			line = w
		case runewidth.StringWidth(line)+1+runewidth.StringWidth(w) > width:
			lines = append(lines, line)
			line = w
		default:
			line += " " + w
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func bubbleLines(greeting, fact string) ([]string, int) {
	lines := append([]string{strings.TrimSpace(greeting)}, wrap(fact, 40)...)
	w := 0
	for _, l := range lines {
		w = max(w, runewidth.StringWidth(l))
	}
	return lines, w
}

func bubble(greeting, fact string) string {
	return frame(bubbleLines(greeting, fact))
}

func frame(lines []string, w int) string {
	var b strings.Builder
	b.WriteString(" ╭")
	b.WriteString(strings.Repeat("─", w+2))
	b.WriteString("╮\n")
	for _, l := range lines {
		b.WriteString(" │ ")
		b.WriteString(runewidth.FillRight(l, w))
		b.WriteString(" │\n")
	}
	b.WriteString(" ╰")
	b.WriteString(strings.Repeat("─", w+2))
	b.WriteString("╯\n")
	b.WriteString(strings.Join(beanArt, "\n"))
	return b.String()
}

func liveBubble(greeting, fact string, delay time.Duration) {
	lines, w := bubbleLines(greeting, fact)
	fmt.Println(frame(make([]string, len(lines)), w))
	height := len(lines) + 2 + len(beanArt)
	for i, l := range lines {
		up := height - 1 - i
		fmt.Printf("\x1b[%dA\x1b[4G", up)
		for _, r := range l {
			fmt.Print(string(r))
			time.Sleep(delay)
		}
		fmt.Printf("\x1b[%dB\r", up)
	}
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	return (err == nil) && (fi.Mode()&os.ModeCharDevice != 0)
}

func main() {
	live := flag.Bool("live", false, "print fact as stream")
	plain := flag.Bool("plain", false, "print only the fact, no bubble")
	flag.Parse()

	fact, err := fetchFact()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not fetch a fact:", err)
		os.Exit(1)
	}

	if *plain || !isTerminal(os.Stdout) {
		fmt.Println(fact)
		return
	}

	if *live {
		liveBubble(greeting, fact, typeDelay)
		return
	}

	fmt.Println(bubble(greeting, fact))
}
