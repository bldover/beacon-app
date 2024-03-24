package screens

type Screen interface {
	Title() string
	Actions() []string
	NextScreen(int) (Screen, *ScreenContext)
}

type ContextScreen interface {
    Screen
	AddContext(ScreenContext)
}

type ContextType int

const (
	Normal = iota
	Selector
)

type ScreenContext struct {
    ReturnScreen Screen
	Props []any
	ContextType ContextType
}

func NewScreenContext(returnScreen Screen, props ...any) *ScreenContext {
    return &ScreenContext{ReturnScreen: returnScreen, Props: props}
}
