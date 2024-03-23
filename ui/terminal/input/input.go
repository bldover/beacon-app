package input

import (
	"bufio"
	"concert-manager/ui/terminal/output"
	"os"
	"strconv"
)

func PromptAndGetInput(prompt string, isValid func(string) error) string {
	output.Displayf("Enter %s:\n", prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		output.Display(">> ")
		in, err := reader.ReadString('\n')
		if err != nil {
			output.Displayln("Error while reading input, try again:")
			continue
		}

		val := in[:len(in) - 1]
		if err := isValid(val); err != nil {
			output.Displayf("Invalid input: %v, try again:\n", err)
			continue
		}

		return val
	}
}

// Valid numbers are in [lowerLimit, upperLimit)
func PromptAndGetInputNumeric(prompt string, lowerLimit int, upperLimit int) int {
	output.Displayf("Enter %s:\n", prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		output.Display(">> ")
		in, err := reader.ReadString('\n')
		if err != nil {
			output.Displayln("Error while reading input, try again:")
			continue
		}

		val, err := strconv.Atoi(in[:len(in) - 1])
		if err != nil {
			output.Displayln("Invalid option, try again:")
			continue
		}
		if val >= upperLimit || val < lowerLimit {
			output.Displayln("Invalid option, try again:")
			continue
		}

		return val
	}
}
