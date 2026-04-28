package article

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestArticle_ToScreen(t *testing.T) {
	article := New().SetName("MyScreen")
	screen := article.ToScreen()

	screen_test.Helper_ToScreen(t, screen)
}

func TestArticle_Stack(t *testing.T) {
	stack := New().
		ToScreen().
		Stack()

	assert.True(t, stack.Has(name))
}

func TestNewArticle_DefaultValues(t *testing.T) {
	article := New()

	assert.Equal(t, name, article.ToScreen().Name())
	assert.Len(t, 0, article.title)
	assert.Len(t, 0, article.article)
}

func TestArticle_SetName(t *testing.T) {
	article := New()
	result := article.SetName("CustomName")

	assert.Equal(t, "CustomName", article.ToScreen().Name())
	assert.Equal(t, result, article)
}

func TestArticle_AddTitleAndArticle(t *testing.T) {
	title := text.NewLine("Title")
	body := text.NewLine("Body")

	article := New().
		AddTitle(*title).
		AddArticle(*body)

	assert.Len(t, 1, article.title)
	assert.Equal(t, text.LineToString(title), text.LineToString(&article.title[0]))

	assert.Len(t, 1, article.article)
	assert.Equal(t, text.LineToString(body), text.LineToString(&article.article[0]))
}

func TestArticle_View(t *testing.T) {
	title := text.NewLine("Title")
	body := text.NewLine("Body")

	article := New().
		AddTitle(*title).
		AddArticle(*body)

	state := state.NewUIState()

	vm := article.view(*state)

	size := winsize.Winsize{
		Rows: 3,
		Cols: 10,
	}

	header := vm.Header.ToDrawable()

	header.Init()
	headers, _ := header.Draw(size)

	kernel := vm.Kernel.ToDrawable()

	kernel.Init()
	lines, _ := kernel.Draw(size)

	assert.Len(t, 1, headers)
	assert.Equal(t, text.LineToString(title), text.LineToString(&headers[0]))

	assert.Len(t, 1, lines)
	assert.Equal(t, text.LineToString(body), text.LineToString(&lines[0]))
}

func TestArticle_Update(t *testing.T) {
	article := New()
	initialState := &state.UIState{}

	article.update(initialState, screen.ScreenEvent{})

	assert.Equal(t, initialState, initialState)
}
