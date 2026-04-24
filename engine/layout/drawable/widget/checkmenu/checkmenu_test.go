package checkmenu

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestCheckMenu_DrawableBasicSuite(t *testing.T) {
	dw := CheckMenuDrawableOptions([]input.CheckOption{})
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
