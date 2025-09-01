package screens

type UtilMenu struct {
	PassedEventManager Screen
	actions            []string
}

const (
	passedEvents = iota + 1
	utilToMainMenu
)

func NewUtilMenu() *UtilMenu {
	menu := UtilMenu{}
	menu.actions = []string{"Manage Passed Events", "Main Menu"}
	return &menu
}

func (m UtilMenu) Title() string {
	return "Utilities"
}

func (m UtilMenu) Actions() []string {
	return m.actions
}

func (m UtilMenu) NextScreen(i int) Screen {
	switch i {
	case passedEvents:
		return m.PassedEventManager
	case utilToMainMenu:
		return nil
	}
	return m
}
