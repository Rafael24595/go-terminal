package help

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/help"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestHelp_DrawableBasicSuite(t *testing.T) {
	dw := HelpDrawableFromMeta(help.NewHelpMeta())
	drawable_test.Test_DrawableBasicSuite(t, dw)
}
