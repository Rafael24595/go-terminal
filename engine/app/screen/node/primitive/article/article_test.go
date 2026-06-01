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

func TestArticle_ToNode(t *testing.T) {
	node := New().Name("base").ToNode()

	screen_test.Helper_ToNode(t, node)
}

func TestArticle_Stack(t *testing.T) {
	stack := New().ToNode().Stack

	assert.True(t, stack.Has(Name))
}

func TestNewArticle_DefaultValues(t *testing.T) {
	article := New()

	assert.Equal(t, Name, article.ToNode().Name)
	assert.Len(t, 0, article.article)
}

func TestArticle_SetName(t *testing.T) {
	article := New()
	result := article.Name("CustomName")

	assert.Equal(t, "CustomName", article.ToNode().Name)
	assert.Equal(t, result, article)
}

func TestArticle_AddTitleAndArticle(t *testing.T) {
	body := text.NewLine("Body")

	article := New().
		AddArticle(*body)

	assert.Len(t, 1, article.article)
	assert.Equal(t, text.LineToString(body), text.LineToString(&article.article[0]))
}

func TestArticle_View(t *testing.T) {
	body := text.NewLine("Body")

	article := New().
		AddArticle(*body)

	state := state.NewUIState()

	vm := article.view(*state)

	size := winsize.Winsize{
		Rows: 3,
		Cols: 10,
	}

	header := vm.Header.ToUnit()

	header.Drawable.Init()

	kernel := vm.Kernel.ToUnit()

	kernel.Drawable.Init()
	lines, _ := kernel.Drawable.Draw(size)

	assert.Len(t, 1, lines)
	assert.Equal(t, text.LineToString(body), text.LineToString(&lines[0]))
}

func TestArticle_Tick(t *testing.T) {
	article := New()
	initialState := &state.UIState{}

	article.tick(initialState, screen.Event{})

	assert.Equal(t, initialState, initialState)
}
