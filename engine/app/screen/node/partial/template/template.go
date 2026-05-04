package template

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const Name = "template"

type Template struct {
	reference string
	stack     set.Set[string]
	model     viewmodel.ViewModel
}

func New() *Template {
	return &Template{
		reference: Name,
		stack:     set.NewSet[string](),
		model:     viewmodel.ViewModel{},
	}
}

func (c *Template) Name(reference string) *Template {
	c.reference = reference
	return c
}

func (c *Template) Stack(stack set.Set[string]) *Template {
	c.stack = stack
	return c
}

func (c *Template) ViewModel(model viewmodel.ViewModel) *Template {
	c.model = model
	return c
}

func (c *Template) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.reference).
		NameToStack().
		AddStack(c.stack).
		Definition(c.definition).
		Update(c.update).
		View(c.view).
		ToNode()
}

func (c *Template) definition() screen.Definition {
	return screen.DefinitionFromKeys()
}

func (c *Template) update(stt *state.UIState, _ screen.ScreenEvent) screen.Result {
	return screen.ResultFromUIState(stt)
}

func (c *Template) view(stt state.UIState) viewmodel.ViewModel {
	return *c.model.Clone()
}
