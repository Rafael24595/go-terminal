package input

import (
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

type MenuOption struct {
	Label  text.Fragment
	Action func() screen.Screen
}

func NewMenuOption(option text.Fragment, action func() screen.Screen) MenuOption {
	return MenuOption{
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
