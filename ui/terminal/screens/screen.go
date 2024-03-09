package screens

type Screen interface {
	Title() string
	Actions() []string
	NextScreen(int) Screen
}
