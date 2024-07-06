package screens

const PageSize = 10

type SortType int

const (
	dateAsc = iota
	dateDesc
)

type Screen interface {
	Title() string
	Actions() []string
	NextScreen(int) Screen
}
