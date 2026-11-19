package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type factMsg struct {
	fact string
	err  error
}

type typeMsg struct{ gen int }

type nextMsg struct{ gen int }

type model struct {
	fetch         func() (string, error)
	interval      time.Duration
	width, height int
	fact          string
	shown         int
	gen           int
}

func (m model) fetchCmd() tea.Msg {
	fact, err := m.fetch()
	return factMsg{fact, err}
}

func typeCmd(gen int) tea.Cmd {
	return tea.Tick(typeDelay, func(time.Time) tea.Msg { return typeMsg{gen} })
}

func nextCmd(gen int, d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg { return nextMsg{gen} })
}

func (m model) Init() tea.Cmd {
	return m.fetchCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case " ", "n":
			return m, m.fetchCmd
		}
	case factMsg:
		m.gen++
		if msg.err != nil {
			return m, nextCmd(m.gen, m.interval)
		}
		m.fact, m.shown = msg.fact, 0
		return m, typeCmd(m.gen)
	case typeMsg:
		if msg.gen != m.gen {
			return m, nil
		}
		m.shown++
		if m.shown >= m.total() {
			return m, nextCmd(m.gen, m.interval)
		}
		return m, typeCmd(m.gen)
	case nextMsg:
		if msg.gen == m.gen {
			return m, m.fetchCmd
		}
	}
	return m, nil
}

func (m model) total() int {
	lines, _ := bubbleLines(greeting, m.fact)
	return runeCount(lines)
}

func (m model) View() string {
	block := strings.Join(beanArt, "\n")
	if m.fact != "" {
		lines, w := bubbleLines(greeting, m.fact)
		block = frame(reveal(lines, m.shown), w)
	}
	return place(block, m.width, m.height)
}

func runTUI(fetch func() (string, error), interval time.Duration) error {
	_, err := tea.NewProgram(model{fetch: fetch, interval: interval}, tea.WithAltScreen()).Run()
	return err
}
