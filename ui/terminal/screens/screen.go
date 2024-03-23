package screens

type Screen interface {
	Title() string
	Actions() []string
	NextScreen(int) Screen
}

type ContextScreen interface {
    Screen
	AddContext(Screen, ...any)
}
