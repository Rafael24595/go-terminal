package line

import (
	"testing"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestLine_DrawableBasicSuite(t *testing.T) {
	dw := LineDrawableFromLines()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
