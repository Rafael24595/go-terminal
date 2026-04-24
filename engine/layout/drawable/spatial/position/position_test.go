package position

import (
	"testing"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestPosition_DrawableBasicSuite(t *testing.T) {
	mock := &drawable_test.MockDrawable{}
	dw := PositionDrawableFromDrawable(mock.ToDrawable())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
