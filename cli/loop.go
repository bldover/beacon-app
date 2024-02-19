package cli

import (
	"concert-manager/out"
	"strings"
)

type Screen interface {
	Title() string
	Actions() []string
	NextScreen(int) Screen
}

type dataDisplayer interface {
    DisplayData()
}

func RunCLI(start Screen) {
	curr := start
    for {
		out.Displayln("----------------------------------------------------------------------")
		out.Displayln(strings.ToUpper(curr.Title()))

		if displayer, ok := curr.(dataDisplayer); ok {
			displayer.DisplayData()
		}

		out.Displayln("Options:")
		actions := curr.Actions()
		spacing := " "
		for i, action := range actions {
			if i >= 9 {
				spacing = ""
			}
			out.Displayf("%s[%d] %s\n", spacing, i + 1, action)
		}

		in := PromptAndGetInputNumeric("option index", 1, len(actions) + 1)
		curr = curr.NextScreen(in)
	}
}
