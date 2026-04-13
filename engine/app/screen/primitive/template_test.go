package primitive

import (
	"testing"

	screen_test "github.com/Rafael24595/go-terminal/test/engine/app/screen"
)

func TestTemplate_ToScreen(t *testing.T) {
	article := NewTemplateScreen()
	screen := article.ToScreen()

	screen_test.Helper_ToScreen(t, screen)
}
