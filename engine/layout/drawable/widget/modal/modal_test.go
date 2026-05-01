package modal

import (
	"testing"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestModal_DrawableBasicSuite(t *testing.T) {
	dw := New().ToDrawable()
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
