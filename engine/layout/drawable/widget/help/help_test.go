package help

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/key"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestHelp_UnitBasicSuite(t *testing.T) {
	unit := UnitFromFields([]key.Descriptor{})
	drawable_test.Test_UnitBasicSuite(t, unit)
}
