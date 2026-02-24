package drawable_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func Helper_ToDrawable(t *testing.T, drawable core.Drawable) {
	t.Helper()

	assert.NotNil(t, drawable.Init, "Drawable.Init should be set")
	assert.NotNil(t, drawable.Draw, "Drawable.Draw should be set")
}
