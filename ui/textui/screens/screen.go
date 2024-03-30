package screens

const pageSize = 10

type sortType int

const (
	dateAsc = iota
	dateDesc
)

type Screen interface {
	Title() string
	Actions() []string
	NextScreen(int) Screen
}
