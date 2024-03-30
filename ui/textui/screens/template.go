package screens

import "concert-manager/ui/textui/output"

type Template struct {
	nextScreen             Screen
	actions                []string
}

const (
	action = iota + 1
	back
)

func NewTemplate() *Template {
	template := Template{}
	template.actions = []string{}
	return &template
}

func (t Template) Title() string {
	return "Template"
}

// optional
func (t Template) DisplayData() {
	output.Displayln("Data")
}

func (t Template) Actions() []string {
	return t.actions
}

func (t Template) NextScreen(i int) Screen {
	switch i {
	case action:
		return t.nextScreen
	case back:
		return nil
	}
	return t
}
