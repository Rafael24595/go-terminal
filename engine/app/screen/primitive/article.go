package primitive

import (
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/render/text"
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

func (c *Article) update(state *state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
	return screen.ScreenResultFromUIState(state)
}

func (c *Article) view(state state.UIState) viewmodel.ViewModel {
	vm := viewmodel.ViewModelFromUIState(state)

	vm.Header.Shift(
		line.EagerDrawableFromLines(c.title...),
	)
	vm.Kernel.Shift(
		line.LazyDrawableFromLines(c.article...),
	)

	return *vm
}
