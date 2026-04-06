package static

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/layout/drawable/block"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestStatic_DrawableBasicSuite(t *testing.T) {
	dw := StaticDrawableFromDrawable(block.BlockDrawableFromLines())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
