package checkmenu

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/model/input"
	
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestCheckMenu_ToDrawable(t *testing.T) {
	dw := CheckMenuDrawableOptions([]input.CheckOption{})
	drawable_test.Helper_ToDrawable(t, dw)
}

func TestCheckMenuDrawable_Draw_ShouldPanicIfNotInitialized(t *testing.T) {
	bd := NewCheckMenuDrawable([]input.CheckOption{})

	assert.Panic(t, func() {
		bd.draw()
	})
}
