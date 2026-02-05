package wrapper_render

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestRenderLineFragments_MergeSameStyles(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("G", core.Bold),
		core.NewFragment("o", core.Bold),
		core.NewFragment("l", core.Bold),
		core.NewFragment("a", core.Bold),
		core.NewFragment("n", core.Bold),
		core.NewFragment("g", core.Bold),
	)

	out := renderLineFragments(line)

	assert.Equal(t, applystyles("Golang", core.Bold), out)
}

func TestRenderLineFragments_StyleChange(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("Hi", core.Bold),
		core.NewFragment(" "),
		core.NewFragment("Ziglang", core.Select),
	)

	out := renderLineFragments(line)

	expected :=
		applystyles("Hi", core.Bold) +
			applystyles(" ") +
			applystyles("Ziglang", core.Select)

	assert.Equal(t, expected, out)
}

func TestRenderLineFragments_DoNotMergeNonContiguous(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("R", core.Bold),
		core.NewFragment("us"),
		core.NewFragment("t", core.Bold),
	)

	out := renderLineFragments(line)

	expected :=
		applystyles("R", core.Bold) +
			applystyles("us") +
			applystyles("t", core.Bold)

	assert.Equal(t, expected, out)
}

func TestRenderLineFragments_EmptyLine(t *testing.T) {
	line := core.Line{}

	out := renderLineFragments(line)

	assert.Equal(t, "", out)
}
