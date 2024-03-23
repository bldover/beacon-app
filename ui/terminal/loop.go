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
	var context *screens.ScreenContext
    for {
		screenChange := curr.Title() != last.Title()
		output.Displayln("----------------------------------------------------------------------")
		output.Displayln(strings.ToUpper(curr.Title()))

		if contextScreen, ok := curr.(screens.ContextScreen); ok && screenChange && context != nil {
			contextScreen.AddContext(*context)
		}

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
		last = curr
		curr, context = curr.NextScreen(in)
	}
}
