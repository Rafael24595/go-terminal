package template

import (
	"testing"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTemplate_ToNode(t *testing.T) {
	node := New().ToNode()

	screen_test.Helper_ToNode(t, node)
}
