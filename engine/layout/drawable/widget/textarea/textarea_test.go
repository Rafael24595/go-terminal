package textarea

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/model/input"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestTextArea_DrawableBasicSuite(t *testing.T) {
	dw := TextAreaDrawableFromData([]rune{}, input.NewTextCursor(false))
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
