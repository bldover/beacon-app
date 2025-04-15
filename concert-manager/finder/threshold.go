package finder

import (
	"fmt"
)

type RecLevel string

const (
	NoMinRec = RecLevel("None")
	LowMinRec = RecLevel("Low")
	MediumMinRec = RecLevel("Medium")
	HighMinRec = RecLevel("High")
)

const (
	noThreshold     = 0
	lowThreshold    = 0.03
	mediumThreshold = 0.08
	highThreshold   = 0.15
)

func ToRecLevel(rank float64) RecLevel {
	switch {
	case rank >= highThreshold:
		return HighMinRec
	case rank >= mediumThreshold:
		return MediumMinRec
	case rank >= lowThreshold:
		return LowMinRec
	case rank >= noThreshold:
		return NoMinRec
	default:
		return "Invalid"
	}
}

func ToThreshold(level RecLevel) (float64, error) {
    switch (level) {
    case NoMinRec:
		return noThreshold, nil
	case LowMinRec:
		return lowThreshold, nil
	case MediumMinRec:
		return mediumThreshold, nil
	case HighMinRec:
		return highThreshold, nil
	default:
		return 0, fmt.Errorf("invalid recommendation level %s", level)
	}
}
