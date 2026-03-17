package cleaner_test

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/cleaner"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func Helper_ToStateCleaner(t *testing.T, cleaner cleaner.StateCleaner) {
	t.Helper()

	assert.NotNil(t, cleaner.Cleanup, "StateCleaner.Cleanup")
}
