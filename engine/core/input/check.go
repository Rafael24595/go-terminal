package input

import (
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/core/text"
)

type CheckOption struct {
	Status bool
	Label  text.Fragment
	action func() screen.Screen
}

func FragmentFromCheckOption(options ...CheckOption) []text.Fragment {
	lines := make([]text.Fragment, len(options))
	for i := range options {
		lines[i] = options[i].Label
	}
	return lines
}
