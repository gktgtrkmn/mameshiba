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
	renderer      liveRenderer
	width, height int
	typing        typing
	screen        string
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
		m.typing = newTyping(greeting, msg.fact)
		m.screen = m.renderer.start(m.typing)
		return m, typeCmd(m.gen)
	case typeMsg:
		if msg.gen != m.gen {
			return m, nil
		}
		m.typing = m.typing.next()
		m.screen = m.renderer.step(m.typing)
		if m.typing.done() {
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

func (m model) View() string {
	if m.screen == "" {
		return place(strings.Join(beanArt, "\n"), m.width, m.height)
	}
	return place(m.screen, m.width, m.height)
}

func runTUI(fetch func() (string, error), interval time.Duration) error {
	m := model{fetch: fetch, interval: interval, renderer: tuiLive{}}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}
