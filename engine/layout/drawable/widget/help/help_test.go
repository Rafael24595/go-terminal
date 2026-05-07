package help

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/key"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestHelp_DrawableBasicSuite(t *testing.T) {
	dw := DrawableFromFields([]key.Descriptor{})
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
