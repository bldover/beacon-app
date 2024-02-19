package cli

import (
	"bufio"
	"concert-manager/out"
	"os"
	"strconv"
)

func NoValidation(_ string) bool {
    return true
}

func PromptAndGetInput(prompt string, isValid func(string) bool) string {
	out.Displayf("Enter %s:\n", prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		out.Display(">> ")
		in, err := reader.ReadString('\n')
		if err != nil {
			out.Displayln("Error while reading input, try again:")
			continue
		}

		val := in[:len(in) - 1]
		if isValid != nil && !isValid(val) {
			out.Displayln("Invalid input, try again:")
			continue
		}

		return val
	}
}

// Valid numbers are in [lowerLimit, upperLimit)
func PromptAndGetInputNumeric(prompt string, lowerLimit int, upperLimit int) int {
	out.Displayf("Enter %s:\n", prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		out.Display(">> ")
		in, err := reader.ReadString('\n')
		out.Display(in)
		if err != nil {
			out.Displayln("Error while reading input, try again:")
			continue
		}

		val, err := strconv.Atoi(in[:len(in) - 1])
		if err != nil {
			out.Displayln("Invalid option, try again:")
			continue
		}
		if val >= upperLimit || val < lowerLimit {
			out.Displayln("Invalid option, try again:")
			continue
		}

		return val
	}
}
