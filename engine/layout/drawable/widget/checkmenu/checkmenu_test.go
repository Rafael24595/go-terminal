package checkmenu

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestCheckMenu_UnitBasicSuite(t *testing.T) {
	unit := UnitFromOptions([]input.CheckOption{})
	drawable_test.Test_UnitBasicSuite(t, unit)
}
