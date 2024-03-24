package format

import (
	"fmt"
	"strings"
)

// adds leading zeros if needed
func formatDate(date string) string {
	parts := strings.Split(date, "/")
	month := parts[0]
	day := parts[1]
	year := parts[2]
	if len(month) == 1 {
		month = "0" + month
	}
	if len(day) == 1 {
		day = "0" + day
	}
	return fmt.Sprintf("%s/%s/%s", month, day, year)
}
