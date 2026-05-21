package position

import (
	"testing"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestPosition_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := UnitFromUnit(mock.ToUnit())
	drawable_test.Test_UnitBasicSuite(t, unit)
}
