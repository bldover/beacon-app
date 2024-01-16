package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Screen interface {
	Title() string
	Data() string
	Actions() []string
	NextScreen(int) Screen
	Parent() Screen
}

func RunCLI(start Screen) {
	curr := start
	var prev Screen
    for {
		fmt.Println("----------------------------------------------------------------------")
		if curr != prev {
			prev = curr
			fmt.Println(strings.ToUpper(curr.Title()))
		}

		data := curr.Data()
		if data != "" {
			fmt.Println(curr.Data())
		}

		fmt.Println("Options:")
		actions := curr.Actions()
		actionCount := len(actions)
		spacing := " "
		for i, action := range actions {
			if i == 9 {
				spacing = ""
			}
			fmt.Printf("%s[%d] %s\n", spacing, i + 1, action)
		}

		if spacing == "" {
			spacing = " "
		}
		hasPrev := curr.Parent() != nil
		if hasPrev {
			fmt.Printf("%s[0] Previous Menu", spacing)
		}
		fmt.Println()

		opt := getUserOption(actionCount, hasPrev)
		if opt == 0 {
			curr = curr.Parent()
		} else {
			curr = curr.NextScreen(opt)
		}
	}
}

func getUserOption(numOpts int, hasPrev bool) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("SELECTION >> ")
		in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading selection, try again:")
			continue
		}

		nbr, err := strconv.Atoi(in[:len(in) - 1])
		if err != nil {
			fmt.Println("Invalid option, try again:", err)
			continue
		}
		if nbr == 0 && !hasPrev {
			fmt.Println("Invalid option, try again:")
			continue
		}
		if nbr > numOpts {
			fmt.Println("Invalid option, try again:")
			continue
		}

		return nbr
	}
}
