package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type MenuOptionAction = func() screen.Screen

type MenuOption struct {
	Id     string
	Label  text.Fragment
	Action func() screen.Screen
}

func NewMenuOption(id string, option text.Fragment, action MenuOptionAction) MenuOption {
	return MenuOption{
		Id:     id,
		Label:  option,
		Action: action,
	}
}

func NewMenuOptions(options ...MenuOption) []MenuOption {
	return options
}

func FragmentFromMenuOption(options ...MenuOption) []text.Fragment {
	lines := make([]text.Fragment, len(options))
	for i := range options {
		lines[i] = options[i].Label
	}
	return lines
}
