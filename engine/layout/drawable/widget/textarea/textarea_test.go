package textarea

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestTextArea_UnitBasicSuite(t *testing.T) {
	unit := New([]rune{}, input.NewTextCursor(false)).ToUnit()
	drawable_test.Test_UnitBasicSuite(t, unit)
}
