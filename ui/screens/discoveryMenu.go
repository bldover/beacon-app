package screens

type DiscoveryMenu struct {
	DiscoveryViewScreen      Screen
	RecommendationViewScreen Screen
	actions                  []string
}

const (
	viewAllUpcoming = iota + 1
	viewRecommended
	discoveryMenuToMainMenu
)

func NewDiscoveryMenu() *DiscoveryMenu {
	menu := DiscoveryMenu{}
	menu.actions = []string{"All Upcoming Events", "Recommended Events", "Main Menu"}
	return &menu
}

func (m DiscoveryMenu) Title() string {
	return "Discovery Menu"
}

func (m DiscoveryMenu) Actions() []string {
	return m.actions
}

func (m DiscoveryMenu) NextScreen(i int) Screen {
	switch i {
	case viewAllUpcoming:
		return m.DiscoveryViewScreen
	case viewRecommended:
		return m.RecommendationViewScreen
	case discoveryMenuToMainMenu:
		return nil
	}
	return m
}
