package help

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestHelp_UnitBasicSuite(t *testing.T) {
	unit := UnitFromFields([]key.Descriptor{})
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestHelpUnit_EmptyFields(t *testing.T) {
	unit := New([]key.Descriptor{}).ToUnit()

	unit.Drawable.Init()

	lines, hasNext := unit.Drawable.Draw(
		winsize.New(5, 80),
	)

	assert.False(t, hasNext)
	assert.Len(t, 0, lines)
}

func TestHelpUnit_WithFields(t *testing.T) {
	unit := New([]key.Descriptor{
		{Code: []string{"RET"}, Detail: "New line/Accept"},
		{Code: []string{"←"}, Detail: "Move left"},
		{Code: []string{"M-b", "Alt-b"}, Detail: "Back"},
	}).ToUnit()

	unit.Drawable.Init()

	cols := winsize.Cols(120)
	lines, hasNext := unit.Drawable.Draw(
		winsize.New(10, cols),
	)

	assert.False(t, hasNext)
	assert.Len(t, 3, lines)

	assert.Len(t, 2, lines[0].Text)
	assert.Equal(t, "--Help---", text.LineToString(&lines[0]))

	assert.Len(t, 3, lines[1].Text)
	assert.Equal(t, "[RET] New line/Accept | ", lines[1].Text[0].Text)
	assert.Equal(t, "[←] Move left | ", lines[1].Text[1].Text)
	assert.Equal(t, "[M-b, Alt-b] Back", lines[1].Text[2].Text)

	assert.Len(t, 1, lines[2].Text)
	assert.Equal(t, "-", text.LineToString(&lines[2]))
}
