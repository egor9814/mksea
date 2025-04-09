package tui

import (
	"fmt"
	"io"
	tui_common "mksea/cmd/sea/internal/common"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listItem string

func (i listItem) FilterValue() string {
	return ""
}

type listItemDelegate struct{}

func (d *listItemDelegate) Height() int {
	return 1
}

func (d *listItemDelegate) Spacing() int {
	return 0
}

func (d *listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d *listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("[%d/%d] %s", index+1, len(m.Items()), i)

	if index == m.Index() {
		str = "> " + str + " <"
	} else {
		str = "  " + str
	}

	fmt.Fprint(w, str)
}

type listModel struct {
	isQuit bool
	list   list.Model
}

func newListModel() tea.Model {
	metaInfo := tui_common.MetaInfo()
	items := make([]list.Item, len(metaInfo.Files))
	for i, it := range metaInfo.Files {
		items[i] = listItem(it)
	}
	l := list.New(items, &listItemDelegate{}, 20, 14)
	l.Title = "content of '" + metaInfo.Name + "'"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return &listModel{
		list: l,
	}
}

func (m *listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esq", "q":
			m.isQuit = true
			return m, tea.Quit

		case "enter":
			// TODO: info
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *listModel) View() string {
	if m.isQuit {
		return ""
	}
	return "\n" + m.list.View()
}
