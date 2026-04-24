package cleaner_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner"
)

func Helper_ToStateCleaner(t *testing.T, cleaner cleaner.StateCleaner) {
	t.Helper()

	assert.NotNil(t, cleaner.Cleanup, "StateCleaner.Cleanup")
}
