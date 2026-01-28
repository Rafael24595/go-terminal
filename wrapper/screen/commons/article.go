package wrapper_commons

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
)

const default_article_name = "Article"

type Article struct {
	reference string
	title      []core.Line
	article    []core.Line
}

func NewArticle() *Article {
	return &Article{
		reference: default_article_name,
		title:      make([]core.Line, 0),
		article:    make([]core.Line, 0),
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

func (c *Article) ToScreen() core.Screen {
	return core.Screen{
		Name:   c.name,
		Update: c.update,
		View:   c.view,
	}
}

func (c *Article) name() string {
	return c.reference
}

func (c *Article) update(state state.UIState, event core.ScreenEvent) core.ScreenResult {
	return core.ScreenResultFromState(state)
}

func (c *Article) view(state state.UIState) core.ViewModel {
	return core.ViewModel{
		Lines: append(c.title, c.article...),
	}
}
