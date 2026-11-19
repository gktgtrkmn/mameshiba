package main

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

const greeting = "ねえ、知ってる？\n"

var beanArt = []string{
	"   ╲",
	"  (◕ ᴥ ◕ ᗜ   )",
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

func reveal(lines []string, n int) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		r := []rune(l)
		k := min(n, len(r))
		out[i] = string(r[:k])
		n -= k
	}
	return out
}

func runeCount(lines []string) int {
	return utf8.RuneCountInString(strings.Join(lines, ""))
}

func place(block string, width, height int) string {
	lines := strings.Split(block, "\n")
	left := strings.Repeat(" ", max(0, (width-runewidth.StringWidth(beanArt[len(beanArt)-1]))/2))
	top := max(0, height/2-(len(lines)-len(beanArt)))
	padded := make([]string, len(lines))
	for i, l := range lines {
		padded[i] = left + l
	}
	return strings.Repeat("\n", top) + strings.Join(padded, "\n")
}
