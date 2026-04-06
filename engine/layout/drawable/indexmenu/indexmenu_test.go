package indexmenu

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/render/text"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestIndexMenu_DrawableBasicSuite(t *testing.T) {
	dw := TextIndexMenuFromData([]text.Fragment{})
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
