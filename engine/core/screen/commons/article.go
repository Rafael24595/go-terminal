package commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
)

const default_article_name = "Article"

type Article struct {
	reference string
	title     []core.Line
	article   []core.Line
}

func NewArticle() *Article {
	return &Article{
		reference: default_article_name,
		title:     make([]core.Line, 0),
		article:   make([]core.Line, 0),
	}
}

func (c *Article) SetName(name string) *Article {
	c.reference = name
	return c
}

func (c *Article) AddTitle(title ...core.Line) *Article {
	c.title = append(c.title, title...)
	return c
}

func (c *Article) AddArticle(article ...core.Line) *Article {
	c.article = append(c.article, article...)
	return c
}

func (c *Article) ToScreen() screen.Screen {
	return screen.Screen{
		Name:       c.name,
		Definition: c.definition,
		Update:     c.update,
		View:       c.view,
	}
}

func (c *Article) name() string {
	return c.reference
}

func (c *Article) definition() screen.Definition {
	return screen.Definition{}
}

func (c *Article) update(state state.UIState, _ screen.ScreenEvent) screen.ScreenResult {
	return screen.ScreenResultFromUIState(state)
}

func (c *Article) view(state state.UIState) core.ViewModel {
	vm := core.ViewModelFromUIState(state)

	vm.Header.Shift(
		line.EagerDrawableFromLines(c.title...),
	)
	vm.Lines.Shift(
		line.LazyDrawableFromLines(c.article...),
	)

	return *vm
}
