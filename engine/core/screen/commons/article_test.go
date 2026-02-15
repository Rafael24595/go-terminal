package commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestArticle_ToScreen(t *testing.T) {
	article := NewArticle().SetName("MyScreen")
	screen := article.ToScreen()

	Helper_ToScreen(t, screen)
}

func TestNewArticle_DefaultValues(t *testing.T) {
	article := NewArticle()

	assert.Equal(t, article.name(), default_article_name, "default name")
	assert.Equal(t, len(article.title), 0, "title lenght")
	assert.Equal(t, len(article.article), 0, "article lenght")
}

func TestArticle_SetName(t *testing.T) {
	article := NewArticle()
	result := article.SetName("CustomName")

	assert.Equal(t, article.name(), "CustomName", "set name")
	assert.Equal(t, result, article, "SetName should return same instance")
}

func TestArticle_AddTitleAndArticle(t *testing.T) {
	title := core.LineFromString("Title")
	body := core.LineFromString("Body")

	article := NewArticle().
		AddTitle(title).
		AddArticle(body)

	assert.Equal(t, len(article.title), 1, "title lines count")
	assert.Equal(t, article.title[0].String(), title.String(), "title content")
	assert.Equal(t, len(article.article), 1, "article lines count")
	assert.Equal(t, article.article[0].String(), body.String(), "article content")
}

func TestArticle_View(t *testing.T) {
	title := core.LineFromString("Title")
	body := core.LineFromString("Body")

	article := NewArticle().
		AddTitle(title).
		AddArticle(body)

	state := state.NewUIState()

	vm := article.view(*state)

	vm.Header.Init(terminal.Winsize{})
	headers, _ := vm.Header.Draw()

	vm.Lines.Init(terminal.Winsize{})
	lines, _ := vm.Lines.Draw()

	assert.Equal(t, len(headers), 1, "ViewModel header count")
	assert.Equal(t, len(lines), 1, "ViewModel lines count")
	assert.Equal(t, headers[0].String(), title.String(), "first line should be title")
	assert.Equal(t, lines[0].String(), body.String(), "second line should be article")
}

func TestArticle_Update(t *testing.T) {
	article := NewArticle()
	initialState := state.UIState{}

	article.update(initialState, screen.ScreenEvent{})

	assert.Equal(t, initialState, initialState, "Update should not change state")
}
