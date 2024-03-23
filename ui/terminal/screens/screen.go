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

type ScreenContext struct {
    ReturnScreen Screen
	Props []any
}

func NewScreenContext(returnScreen Screen, props ...any) *ScreenContext {
    return &ScreenContext{returnScreen, props}
}
