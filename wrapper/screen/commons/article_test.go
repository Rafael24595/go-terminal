package wrapper_commons

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
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

	assert.Equal(t, len(vm.Lines), 2, "ViewModel lines count")
	assert.Equal(t, vm.Lines[0].String(), title.String(), "first line should be title")
	assert.Equal(t, vm.Lines[1].String(), body.String(), "second line should be article")
}

func TestArticle_Update(t *testing.T) {
	article := NewArticle()
	initialState := state.UIState{}

	result := article.update(initialState, core.ScreenEvent{})

	assert.Equal(t, result.State, initialState, "Update should not change state")
}
