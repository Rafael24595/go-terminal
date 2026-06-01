package article

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "article"

type Article struct {
	reference string
	article   []text.Line
}

func New() *Article {
	return &Article{
		reference: Name,
		article:   make([]text.Line, 0),
	}
}

func (n *Article) Name(name string) *Article {
	n.reference = name
	return n
}

func (n *Article) AddArticle(article ...text.Line) *Article {
	n.article = append(n.article, article...)
	return n
}

func (n *Article) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		WithoutKeys().
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *Article) tick(uiState *state.UIState, _ screen.Event) screen.Result {
	return screen.ResultFromUIState(uiState)
}

func (n *Article) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	vm.Kernel.Push(
		line.UnitFromLines(n.article...),
	)

	return *vm
}
