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

func (n *Template) Name(reference string) *Template {
	n.reference = reference
	return n
}

func (n *Template) Stack(stack set.Set[string]) *Template {
	n.stack = stack
	return n
}

func (n *Template) ViewModel(model viewmodel.ViewModel) *Template {
	n.model = model
	return n
}

func (n *Template) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		AddStack(n.stack).
		Definition(n.definition).
		Update(n.update).
		View(n.view).
		ToNode()
}

func (n *Template) definition() screen.Definition {
	return screen.EmptyDefinition()
}

func (n *Template) update(stt *state.UIState, _ screen.Event) screen.Result {
	return screen.ResultFromUIState(stt)
}

func (n *Template) view(stt state.UIState) viewmodel.ViewModel {
	return *n.model.Clone()
}
