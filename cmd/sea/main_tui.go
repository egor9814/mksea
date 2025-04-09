//go:build bubbletea_tui

package main

import "mksea/cmd/sea/internal/tui"

func main_tui() bool {
	tui.Main()
	return true
}
