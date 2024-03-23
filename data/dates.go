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

func EventSorterDateAsc() func(a, b Event) int {
	return func(a, b Event) int {
		if Timestamp(a.Date).Before(Timestamp(b.Date)) {
			return -1
		} else if Timestamp(a.Date).After(Timestamp(b.Date)) {
			return 1
		} else {
			return 0
		}
	}
}

func EventSorterDateDesc() func(a, b Event) int {
	return func(a, b Event) int {
		if Timestamp(a.Date).Before(Timestamp(b.Date)) {
			return 1
		} else if Timestamp(a.Date).After(Timestamp(b.Date)) {
			return -1
		} else {
			return 0
		}
	}
}

func EventDetailsSorterDateAsc() func(a, b EventDetails) int {
	return func(a, b EventDetails) int {
		return EventSorterDateAsc()(a.Event, b.Event)
	}
}

func EventDetailsSorterDateDesc() func(a, b EventDetails) int {
	return func(a, b EventDetails) int {
		return EventSorterDateDesc()(a.Event, b.Event)
	}
}

func ValidPastDate(date string) bool {
	if !ValidDate(date) {
		return false
	}
	return Timestamp(date).Before(time.Now())
}

func ValidFutureDate(date string) bool {
	if !ValidDate(date) {
		return false
	}
	now := time.Now()
	return Timestamp(date).Equal(now) || Timestamp(date).After(now)
}
