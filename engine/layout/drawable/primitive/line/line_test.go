package line

import (
	"testing"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestLine_UnitBasicSuite(t *testing.T) {
	unit := UnitFromLines()
	drawable_test.Test_UnitBasicSuite(t, unit)
}
