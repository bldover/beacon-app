package data

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// expects format "#/#/# but doesn't check for non-standard date values (like day 32)
// since time.Date handles these with overflow. Consumers of the Event type should regulate
// their own date values if possible overflow is undesired
func ValidDate(date string) bool {
	parts := strings.Split(date, "/")
	if len(parts) != 3 {
		return false
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(parts[2]); err != nil {
		return false
	}
	return true
}

func ValidPastDate(date string) bool {
	if !ValidDate(date) {
		return false
	}
	return time.Now().Equal(Timestamp(date)) || time.Now().After(Timestamp(date))
}

func ValidFutureDate(date string) bool {
	if !ValidDate(date) {
		return false
	}
	return time.Now().Before(Timestamp(date))
}

// format is "mm/dd/yyyy", with leading zeros optional
// expected that the date string has been previously validated to not error when converted to ints
func Timestamp(date string) time.Time {
	parts := strings.Split(date, "/")
    month, _ := strconv.Atoi(parts[0])
	day, _ := strconv.Atoi(parts[1])
	year, _ := strconv.Atoi(parts[2])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func Date(ts time.Time) string {
	day, month, year := ts.Day(), ts.Month(), ts.Year()
	return fmt.Sprintf("%d/%d/%d", month, day, year)
}
