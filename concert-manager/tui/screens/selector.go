package screens

import (
	"concert-manager/tui/output"
)

type Selector[T any] struct {
	ScreenTitle  string
	Formatter    sliceFormatter[T]
	HandleSelect selectAction[T]
	Options      []T
	Next         Screen
}

type sliceFormatter[T any] func([]T) []string
type selectAction[T any] func(T)

func IdentityTransform[T any](x []T) []T {
	return x
}

func (s Selector[_]) Title() string {
	return s.ScreenTitle
}

func (s Selector[_]) DisplayData() {
	if len(s.Options) == 0 {
		output.Displayln("No data found")
	}
}

func (s Selector[_]) Actions() []string {
	actions := []string{}
	pageEvents := s.Options
	actions = append(actions, s.Formatter(pageEvents)...)
	actions = append(actions, "Cancel")
	return actions
}

func (s *Selector[_]) NextScreen(i int) Screen {
	if i != len(s.Options)+1 {
		s.HandleSelect(s.Options[i-1])
		return s.Next
	}
	return nil
}
