package template

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

type Template struct {
	reference string
	stack     set.Set[string]
	model     viewmodel.ViewModel
}

func New() *Template {
	return &Template{}
}

func (c *Template) SetName(reference string) *Template {
	c.reference = reference
	return c
}

func (c *Template) SetStack(stack set.Set[string]) *Template {
	c.stack = stack
	return c
}

func (c *Template) SetViewModel(model viewmodel.ViewModel) *Template {
	c.model = model
	return c
}

func (c *Template) ToScreen() screen.Screen {
	screen := screen.Screen{
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}

	screen = screen.SetName(c.reference)

	if len(c.stack) > 0 {
		screen = screen.SetStack(c.stack)
	} else {
		screen = screen.StackFromName()
	}

	return screen
}

func (c *Template) definition() screen.Definition {
	return screen.DefinitionFromKeys()
}

func (c *Template) update(stt *state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
	return screen.ScreenResultFromUIState(stt)
}

func (c *Template) view(stt state.UIState) viewmodel.ViewModel {
	return *c.model.Clone()
}
