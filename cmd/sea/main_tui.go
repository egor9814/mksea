//go:build tui

package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	grid := tview.NewGrid().
		SetRows(2).
		SetColumns(0).
		SetBorders(true)
	app.SetRoot(grid, true).SetFocus(grid)

	box := tview.NewTextView().SetText("Hello").SetTextAlign(tview.AlignCenter)
	grid.AddItem(box, 0, 0, 1, 1, 0, 0, false)

	label := tview.NewTextView().SetText(fmt.Sprintf("exe size: 0x%016X", exe_size()))
	label.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})
	grid.AddItem(label, 1, 0, 1, 1, 0, 100, false)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
