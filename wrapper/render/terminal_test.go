package wrapper_render

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

func TestRenderLineFragments_MergeSameStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("G").AddAtom(style.AtmBold),
		*text.NewFragment("o").AddAtom(style.AtmBold),
		*text.NewFragment("l").AddAtom(style.AtmBold),
		*text.NewFragment("a").AddAtom(style.AtmBold),
		*text.NewFragment("n").AddAtom(style.AtmBold),
		*text.NewFragment("g").AddAtom(style.AtmBold),
	)

	out := renderLineFragments(*line, winsize.Winsize{})

	assert.Equal(t, applyAtomStyles("Golang", style.AtmBold), out)
}

func TestRenderLineFragments_StyleChange(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("Hi").AddAtom(style.AtmBold),
		*text.NewFragment(" "),
		*text.NewFragment("Ziglang").AddAtom(style.AtmSelect),
	)

	out := renderLineFragments(*line, winsize.Winsize{})

	expected :=
		applyAtomStyles("Hi", style.AtmBold) +
			applyAtomStyles(" ") +
			applyAtomStyles("Ziglang", style.AtmSelect)

	assert.Equal(t, expected, out)
}

func TestRenderLineFragments_DoNotMergeNonContiguous(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("R").AddAtom(style.AtmBold),
		*text.NewFragment("us"),
		*text.NewFragment("t").AddAtom(style.AtmBold),
	)

	out := renderLineFragments(*line, winsize.Winsize{})

	expected :=
		applyAtomStyles("R", style.AtmBold) +
			applyAtomStyles("us") +
			applyAtomStyles("t", style.AtmBold)

	assert.Equal(t, expected, out)
}

func TestRenderLineFragments_EmptyLine(t *testing.T) {
	line := text.Line{}

	out := renderLineFragments(line, winsize.Winsize{})

	assert.Equal(t, "", out)
}
