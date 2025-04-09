package tui

import (
	tui_common "mksea/cmd/sea/internal/common"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type passwordModel struct {
	input    textinput.Model
	attempts int
	isQuit   bool
	listMode bool
}

func newPasswordModel(listMode bool) tea.Model {
	pi := textinput.New()
	pi.Cursor.Style = cursorStyle
	pi.Placeholder = "password"
	pi.EchoMode = textinput.EchoPassword
	pi.EchoCharacter = '•'
	pi.PromptStyle = focusedStyle
	pi.TextStyle = focusedStyle
	pi.Focus()

	return &passwordModel{
		input:    pi,
		attempts: tui_common.PasswordAttempts,
		listMode: listMode,
	}
}

func (m *passwordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *passwordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.isQuit = true
			return m, tea.Quit

		case "enter":
			return m.testPassword()
		}
	}
	cmd := m.updateInput(msg)
	return m, cmd
}

func (m *passwordModel) updateInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m *passwordModel) View() string {
	if m.isQuit {
		return ""
	}
	var b strings.Builder
	b.WriteString("type archive password")
	if m.attempts != tui_common.PasswordAttempts {
		b.WriteString(" (attempts left: ")
		b.WriteString(strconv.Itoa(m.attempts))
		b.WriteByte(')')
	}
	b.WriteString(":\n")
	b.WriteString(m.input.View())
	return b.String()
}

func (m *passwordModel) testPassword() (tea.Model, tea.Cmd) {
	p := []byte(m.input.Value())
	if tui_common.DecodeEncoderKey(p) {
		if m.listMode {
			return newListModel(), nil
		}
		return newUnpackModel(), nil
	} else {
		m.attempts--
		if m.attempts == 0 {
			return m, tea.Quit
		}
		m.input.SetValue("")
	}
	return m, nil
}
