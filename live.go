package main

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

type typing struct {
	lines []string
	w     int
	total int
	n     int
}

func newTyping(greeting, fact string) typing {
	lines, w := bubbleLines(greeting, fact)
	return typing{lines: lines, w: w, total: runeCount(lines)}
}

func (t typing) done() bool {
	return t.n >= t.total
}

func (t typing) next() typing {
	t.n = min(t.n+1, t.total)
	return t
}

func (t typing) height() int {
	return len(t.lines) + 2 + len(beanArt)
}

func (t typing) last() (line, col int, r rune) {
	idx := t.n - 1
	for i, l := range t.lines {
		rs := []rune(l)
		if idx < len(rs) {
			return i, runewidth.StringWidth(string(rs[:idx])), rs[idx]
		}
		idx -= len(rs)
	}
	return -1, 0, 0
}

type liveRenderer interface {
	start(t typing) string
	step(t typing) string
}

type cliLive struct{}

func (cliLive) start(t typing) string {
	return frame(reveal(t.lines, t.n), t.w) + "\n"
}

func (cliLive) step(t typing) string {
	line, col, r := t.last()
	if line < 0 {
		return ""
	}
	up := t.height() - 1 - line
	return fmt.Sprintf("\x1b[%dA\x1b[%dG%c\x1b[%dB\r", up, 4+col, r, up)
}

type tuiLive struct{}

func (tuiLive) start(t typing) string {
	return frame(reveal(t.lines, t.n), t.w)
}

func (l tuiLive) step(t typing) string {
	return l.start(t)
}
