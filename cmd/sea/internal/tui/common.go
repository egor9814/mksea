package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Dark:  "#e0e0e0",
		Light: "#292929",
	})
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	itemStyle           = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.AdaptiveColor{
		Dark:  "#29e029",
		Light: "#e029e0",
	})
	// titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	titleStyle      = noStyle
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)
