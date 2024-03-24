package output

import (
	"concert-manager/log"
	"fmt"
)

func Display(v ...any) {
	log.Display(v...)
	fmt.Print(v...)
}

func Displayf(format string, v ...any) {
	log.Displayf(format, v...)
	fmt.Printf(format, v...)
}

func Displayln(v ...any) {
	log.Display(v...)
	fmt.Println(v...)
}

func ClearCurrentLine() {
	log.Display("Cleared last line")
	fmt.Printf("\r\033[K") // carriage return + kill line, probably platform-dependent
}
