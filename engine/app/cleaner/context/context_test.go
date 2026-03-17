package context

import (
	"testing"

	cleaner_test "github.com/Rafael24595/go-terminal/test/engine/app/cleaner"
)

func TestContext_ToStateCleaner(t *testing.T) {
	cleaner := NewContextCleaner()

	cleaner_test.Helper_ToStateCleaner(t, cleaner)
}
