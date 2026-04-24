package modal

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestModal_DrawableBasicSuite(t *testing.T) {
	dw := ModalDrawableFromData([]text.Line{}, []text.Fragment{}, 0)
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
