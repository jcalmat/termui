/// +build ignore

package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type nodeValue string

func (nv nodeValue) String() string {
	return string(nv)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	// defer ui.Close()

	cava := widgets.NewCheckbox("Ça va ?", false)
	bien := widgets.NewCheckbox("Bien ?", false)
	bienbien := widgets.NewCheckbox("Bien bien ?", false)
	nodes := []*widgets.FormNode{
		{
			Item: cava,
			Nodes: []*widgets.FormNode{
				{
					Item: bien,
					Nodes: []*widgets.FormNode{
						{
							Item: bienbien,
						},
					},
				},
			},
		},
	}

	l := widgets.NewForm()
	l.Title = "test"
	l.SetNodes(nodes)

	x, y := ui.TerminalDimensions()
	l.SetRect(0, 0, x, y)

	// l.TextStyle = ui.NewStyle(ui.ColorYellow)
	// ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	l.SelectedTextStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	// l.WrapText = false
	// l.SetNodes(nodes)

	// x, y := ui.TerminalDimensions()

	// l.SetRect(0, 0, x, y)

	ui.Render(l)

	// previousKey := ""
	var close bool
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			ui.Close()
			close = true
			break
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		// case "<C-d>":
		// 	l.ScrollHalfPageDown()
		// case "<C-u>":
		// 	l.ScrollHalfPageUp()
		// case "<C-f>":
		// 	l.ScrollPageDown()
		// case "<C-b>":
		// 	l.ScrollPageUp()
		// // case "g":
		// // 	if previousKey == "g" {
		// // 		l.ScrollTop()
		// // 	}
		// case "<Home>":
		// 	l.ScrollTop()
		case "<Enter>":
			l.ToggleExpand()
			// case "G", "<End>":
			// 	l.ScrollBottom()
			// case "E":
		// 	l.ExpandAll()
		// case "C":
		// 	l.CollapseAll()
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			l.SetRect(0, 0, x, y)
		}

		// if previousKey == "g" {
		// 	previousKey = ""
		// } else {
		// 	previousKey = e.ID
		// }

		if close {
			break
		}
		ui.Render(l)
	}

	fmt.Printf("cava = %v | bien = %v | bienbien = %v\n", cava.Answer(), bien.Answer(), bienbien.Answer())

}
