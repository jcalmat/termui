/// +build ignore

package main

import (
	"fmt"
	"log"

	ui "github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	var close bool
	checkbox0 := widgets.NewCheckbox("checkbox 0", false)
	checkbox01 := widgets.NewCheckbox("checkbox 0.1", false)
	checkbox02 := widgets.NewCheckbox("checkbox 0.2", false)
	textfield1 := widgets.NewTextField("textfield 1:")
	checkbox11 := widgets.NewCheckbox("checkbox 1.1", false)
	label0 := widgets.NewLabel("label 0")
	label01 := widgets.NewLabel("label 0.1")
	button0 := widgets.NewButton("Close button 0", func() {
		close = true
	})
	nodes := []*widgets.FormNode{
		{
			Item: label0,
		},
		{
			Item: checkbox0,
			Nodes: []*widgets.FormNode{
				{
					Item: checkbox01,
					Nodes: []*widgets.FormNode{
						{
							Item: checkbox02,
						},
						{
							Item: label01,
						},
					},
				},
			},
		},
		{
			Item: textfield1,
			Nodes: []*widgets.FormNode{
				{
					Item: checkbox11,
				},
			},
		},
		{
			Item: button0,
		},
	}

	x, y := ui.TerminalDimensions()

	l := widgets.NewForm()
	l.Title = "My form"
	l.SetNodes(nodes)
	l.SetRect(0, 0, x, y)
	l.SelectedTextStyle = ui.NewStyle(ui.ColorClear)

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			close = true
			break
		case "<Down>":
			l.ScrollDown()
		case "<Up>":
			l.ScrollUp()
		case "<Enter>":
			l.ToggleExpand()
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			l.SetRect(0, 12, x, y)
		}

		l.HandleKeyboard(e)

		if close {
			break
		}
		ui.Render(l)
	}

	ui.Close()

	fmt.Printf("checkbox0 = %v | checkbox01 = %v | checkbox02 = %v | checkbox11 = %v\ntextfield = %s\n", checkbox0.Answer(), checkbox01.Answer(), checkbox02.Answer(), checkbox11.Answer(), textfield1.Answer())

}
