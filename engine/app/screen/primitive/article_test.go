package primitive

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-terminal/engine/app/screen"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestArticle_ToScreen(t *testing.T) {
	article := NewArticle().SetName("MyScreen")
	screen := article.ToScreen()

	screen_test.Helper_ToScreen(t, screen)
}

func TestArticle_Stack(t *testing.T) {
	stack := NewArticle().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(default_article_name))
}

func TestNewArticle_DefaultValues(t *testing.T) {
	article := NewArticle()

	assert.Equal(t, article.ToScreen().Name(), default_article_name, "default name")
	assert.Equal(t, len(article.title), 0, "title lenght")
	assert.Equal(t, len(article.article), 0, "article lenght")
}

func TestArticle_SetName(t *testing.T) {
	article := NewArticle()
	result := article.SetName("CustomName")

	assert.Equal(t, article.ToScreen().Name(), "CustomName", "set name")
	assert.Equal(t, result, article, "SetName should return same instance")
}

func TestArticle_AddTitleAndArticle(t *testing.T) {
	title := text.LineFromString("Title")
	body := text.LineFromString("Body")

	article := NewArticle().
		AddTitle(title).
		AddArticle(body)

	assert.Equal(t, len(article.title), 1, "title lines count")
	assert.Equal(t, text.LineToString(article.title[0]), text.LineToString(title), "title content")
	assert.Equal(t, len(article.article), 1, "article lines count")
	assert.Equal(t, text.LineToString(article.article[0]), text.LineToString(body), "article content")
}

func TestArticle_View(t *testing.T) {
	title := text.LineFromString("Title")
	body := text.LineFromString("Body")

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
	assert.Equal(t, text.LineToString(headers[0]), text.LineToString(title), "first line should be title")
	assert.Equal(t, text.LineToString(lines[0]), text.LineToString(body), "second line should be article")
}

func TestArticle_Update(t *testing.T) {
	article := NewArticle()
	initialState := &state.UIState{}

	article.update(initialState, screen.ScreenEvent{})

	assert.Equal(t, initialState, initialState, "Update should not change state")
}
