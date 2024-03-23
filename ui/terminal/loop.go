package terminal

import (
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"strings"
)

type dataDisplayer interface {
    DisplayData()
}

type refresher interface {
    Refresh()
}

func RunUI(start screens.Screen) {
	curr := start
	last := start
    for {
		output.Displayln("----------------------------------------------------------------------")
		output.Displayln(strings.ToUpper(curr.Title()))

		if refresher, ok := curr.(refresher); ok && curr.Title() != last.Title() {
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
	}
}
