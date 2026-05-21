package box

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestBox_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := UnitFromUnit(mock.ToUnit())
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestBox_Init_ShouldPropagateToChild(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit())

	unit.init()

	assert.True(t, unit.loaded)
	assert.True(t, mock.InitCalled)
}
