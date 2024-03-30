package core

import (
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
	"strings"
)

type dataDisplayer interface {
    DisplayData()
}

type refresher interface {
    Refresh()
}

func Run(start screens.Screen) {
	history := history{[]screens.Screen{start}}
	curr := start
	last := start
	var screenChange bool
    for {
		output.Displayln("----------------------------------------------------------------------")
		output.Displayln(strings.ToUpper(curr.Title()))

		if refresher, ok := curr.(refresher); ok && screenChange {
			refresher.Refresh()
		}

		if displayer, ok := curr.(dataDisplayer); ok {
			displayer.DisplayData()
		}

		output.Displayln("Options:")
		actions := curr.Actions()
		spacing := " "
		for i, action := range actions {
			if i >= 9 {
				spacing = ""
			}
			output.Displayf("%s[%d] %s\n", spacing, i + 1, action)
		}

		in := input.PromptAndGetInputNumeric("option index", 1, len(actions) + 1)
		last, curr = curr, curr.NextScreen(in)

		screenChange = curr == nil || curr.Title() != last.Title()
		if screenChange {
			if curr == nil {
				curr = history.getPrevious()
			} else {
				history.update(curr)
			}
		}
	}
}
