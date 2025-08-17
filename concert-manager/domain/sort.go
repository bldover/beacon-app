package domain

import (
	"concert-manager/util"
)

func EventSorterDateAsc() func(a, b Event) int {
	return func(a, b Event) int {
		if util.Timestamp(a.Date).Before(util.Timestamp(b.Date)) {
			return -1
		} else if util.Timestamp(a.Date).After(util.Timestamp(b.Date)) {
			return 1
		} else {
			return 0
		}
	}
}

func EventSorterDateDesc() func(a, b Event) int {
	return func(a, b Event) int {
		if util.Timestamp(a.Date).Before(util.Timestamp(b.Date)) {
			return 1
		} else if util.Timestamp(a.Date).After(util.Timestamp(b.Date)) {
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
