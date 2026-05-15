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

func (c *Article) Name(name string) *Article {
	c.reference = name
	return c
}

func (c *Article) AddArticle(article ...text.Line) *Article {
	c.article = append(c.article, article...)
	return c
}

func (c *Article) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(c.reference).
		NameToStack().
		WithoutDefinition().
		Update(c.update).
		View(c.view).
		ToNode()
}

func (c *Article) update(stt *state.UIState, _ screen.Event) screen.Result {
	return screen.ResultFromUIState(stt)
}

func (c *Article) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	vm.Kernel.Push(
		line.DrawableFromLines(c.article...),
	)

	return *vm
}
