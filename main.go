package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

const (
	factURL  = "https://uselessfacts.jsph.pl/api/v2/facts/random?language=en"
	greeting = "ねえ、知ってる？\n"
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

func bubble(greeting, fact string) string {
	lines := append([]string{strings.TrimSpace(greeting)}, wrap(fact, 40)...)
	w := 0
	for _, l := range lines {
		w = max(w, runewidth.StringWidth(l))
	}
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

func main() {
	fact, err := fetchFact()
	if err != nil {
		return
	}
	fmt.Println(bubble(greeting, fact))
}
