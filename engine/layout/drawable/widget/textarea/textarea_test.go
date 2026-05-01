package textarea

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestTextArea_DrawableBasicSuite(t *testing.T) {
	dw := New([]rune{}, input.NewTextCursor(false)).ToDrawable()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
