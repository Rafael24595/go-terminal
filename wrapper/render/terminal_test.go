package wrapper_render

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestRenderLineFragments_MergeSameStyles(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("G").AddAtom(style.AtmBold),
		core.NewFragment("o").AddAtom(style.AtmBold),
		core.NewFragment("l").AddAtom(style.AtmBold),
		core.NewFragment("a").AddAtom(style.AtmBold),
		core.NewFragment("n").AddAtom(style.AtmBold),
		core.NewFragment("g").AddAtom(style.AtmBold),
	)

	out := renderLineFragments(line, terminal.Winsize{})

	assert.Equal(t, applyAtomStyles("Golang", style.AtmBold), out)
}

func TestRenderLineFragments_StyleChange(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("Hi").AddAtom(style.AtmBold),
		core.NewFragment(" "),
		core.NewFragment("Ziglang").AddAtom(style.AtmSelect),
	)

	out := renderLineFragments(line, terminal.Winsize{})

	expected :=
		applyAtomStyles("Hi", style.AtmBold) +
			applyAtomStyles(" ") +
			applyAtomStyles("Ziglang", style.AtmSelect)

	assert.Equal(t, expected, out)
}

func TestRenderLineFragments_DoNotMergeNonContiguous(t *testing.T) {
	line := core.LineFromFragments(
		core.NewFragment("R").AddAtom(style.AtmBold),
		core.NewFragment("us"),
		core.NewFragment("t").AddAtom(style.AtmBold),
	)

	out := renderLineFragments(line, terminal.Winsize{})

	expected :=
		applyAtomStyles("R", style.AtmBold) +
			applyAtomStyles("us") +
			applyAtomStyles("t", style.AtmBold)

	assert.Equal(t, expected, out)
}

func TestRenderLineFragments_EmptyLine(t *testing.T) {
	line := core.Line{}

	out := renderLineFragments(line, terminal.Winsize{})

	assert.Equal(t, "", out)
}
