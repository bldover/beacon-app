package util

import "concert-manager/data"

func EventSorterDateAsc() func(a, b data.Event) int {
	return func(a, b data.Event) int {
		if Timestamp(a.Date).Before(Timestamp(b.Date)) {
			return -1
		} else if Timestamp(a.Date).After(Timestamp(b.Date)) {
			return 1
		} else {
			return 0
		}
	}
}

func EventSorterDateDesc() func(a, b data.Event) int {
	return func(a, b data.Event) int {
		if Timestamp(a.Date).Before(Timestamp(b.Date)) {
			return 1
		} else if Timestamp(a.Date).After(Timestamp(b.Date)) {
			return -1
		} else {
			return 0
		}
	}
}

func EventDetailsSorterDateAsc() func(a, b data.EventDetails) int {
	return func(a, b data.EventDetails) int {
		return EventSorterDateAsc()(a.Event, b.Event)
	}
}

func EventDetailsSorterDateDesc() func(a, b data.EventDetails) int {
	return func(a, b data.EventDetails) int {
		return EventSorterDateDesc()(a.Event, b.Event)
	}
}
