package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
)

type TemplateScreen struct {
	reference string
	stack     set.Set[string]
	model     viewmodel.ViewModel
}

func NewTemplateScreen() *TemplateScreen {
	return &TemplateScreen{}
}

func (c *TemplateScreen) SetName(reference string) *TemplateScreen {
	c.reference = reference
	return c
}

func (c *TemplateScreen) SetStack(stack set.Set[string]) *TemplateScreen {
	c.stack = stack
	return c
}

func (c *TemplateScreen) SetViewModel(model viewmodel.ViewModel) *TemplateScreen {
	c.model = model
	return c
}

func (c *TemplateScreen) ToScreen() screen.Screen {
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

func (c *TemplateScreen) definition() screen.Definition {
	return screen.DefinitionFromKeys()
}

func (c *TemplateScreen) update(stt *state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
	return screen.ScreenResultFromUIState(stt)
}

func (c *TemplateScreen) view(stt state.UIState) viewmodel.ViewModel {
	return *c.model.Clone()
}
