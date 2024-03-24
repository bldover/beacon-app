package common

import (
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
)

type Selector[T any, U any] struct {
	ScreenTitle string
	Formatter         sliceFormatter[T]
	OutputTransformer transformer[T, U]
	Next              screens.Screen
	options           []T
	returnScreen      screens.Screen
}

type sliceFormatter[T any] func([]T) []string
type transformer[T any, U any] func(T) U

func IdentityTransform[T any, U any](x T) U {
    return any(x).(U)
}

func (s *Selector[T, _]) AddContext(context screens.ScreenContext) {
	s.returnScreen = context.ReturnScreen
	s.options = context.Props[0].([]T)
}

func (s Selector[_, _]) Title() string {
	return s.ScreenTitle
}

func (s Selector[_, _]) DisplayData() {
	if len(s.options) == 0 {
		output.Displayln("No data found")
	}
}

func (s Selector[_, _]) Actions() []string {
	actions := []string{}
	pageEvents := s.options
	actions = append(actions, s.Formatter(pageEvents)...)
	actions = append(actions, "Cancel")
	return actions
}

func (s *Selector[_, _]) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	if i != len(s.options)+1 {
		out := s.OutputTransformer(s.options[i-1])
		context := screens.NewScreenContext(s.returnScreen, out)
		context.ContextType = screens.Selector
		return s.Next, context
	}
	return s.returnScreen, nil
}
