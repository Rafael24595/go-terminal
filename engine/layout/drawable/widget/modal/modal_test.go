package modal

import (
	"testing"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestModal_UnitBasicSuite(t *testing.T) {
	unit := New().ToUnit()
	drawable_test.Test_UnitBasicSuite(t, unit)
}
