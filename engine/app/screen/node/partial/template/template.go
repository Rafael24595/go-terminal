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
		WithoutKeys().
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *Template) tick(uiState *state.UIState, _ screen.Event) screen.Result {
	return screen.ResultFromUIState(uiState)
}

func (n *Template) view(state.UIState) viewmodel.ViewModel {
	return *n.model.Clone()
}
