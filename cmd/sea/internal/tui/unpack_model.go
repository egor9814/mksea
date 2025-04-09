package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type unpackModel struct {
	isQuit bool
}

func newUnpackModel() tea.Model {
	return &unpackModel{}
}

func (m *unpackModel) Init() tea.Cmd {
	return nil
}

func (m *unpackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.isQuit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *unpackModel) View() string {
	if m.isQuit {
		return ""
	}
	return "hello"
}
