package template

import (
	"testing"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTemplate_ToScreen(t *testing.T) {
	article := New()
	screen := article.ToScreen()

	screen_test.Helper_ToScreen(t, screen)
}
