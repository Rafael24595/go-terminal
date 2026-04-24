package primitive

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const default_article_name = "Article"

type Article struct {
	reference string
	title     []text.Line
	article   []text.Line
}

func NewArticle() *Article {
	return &Article{
		reference: default_article_name,
		title:     make([]text.Line, 0),
		article:   make([]text.Line, 0),
	}
}

func (c *Article) SetName(name string) *Article {
	c.reference = name
	return c
}

func (c *Article) AddTitle(title ...text.Line) *Article {
	c.title = append(c.title, title...)
	return c
}

func (c *Article) AddArticle(article ...text.Line) *Article {
	c.article = append(c.article, article...)
	return c
}

func (c *Article) ToScreen() screen.Screen {
	screen := screen.Screen{
		Update: c.update,
		View:   c.view,
	}

	return screen.SetName(c.reference).
		SetDefinition().
		StackFromName()
}

func (c *Article) update(stt *state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
	return screen.ScreenResultFromUIState(stt)
}

func (c *Article) view(_ state.UIState) viewmodel.ViewModel {
	vm := viewmodel.NewViewModel()

	vm.Header.Push(
		block.BlockDrawableFromLines(c.title...),
	)
	vm.Kernel.Push(
		line.LineDrawableFromLines(c.article...),
	)

	return *vm
}
