package justify

import (
	"strings"
	"testing"

	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func fragmentTexts(frags []text.Fragment) string {
	var s strings.Builder
	for _, f := range frags {
		s.WriteString(f.Text)
	}
	return s.String()
}

func renderFragments(frags []text.Fragment) string {
	var s strings.Builder
	for _, f := range frags {
		s.WriteString(f.Text)

		count := uint(0)
		ok := false
		if args, ex := f.Spec.Args()[style.KeyRepeatRightSize]; ex {
			ok = true
			count = args.Uintd(0)
		}

		if args, ex := f.Spec.Args()[style.KeyPaddingRightSize]; ex {
			ok = true
			count = args.Uintd(0)
		}

		if ok {
			for range count {
				s.WriteString(" ")
			}
		}
	}

	return s.String()
}

func renderLine(cols int, mode style.Justify, line text.Line) string {
	frags := renderFragments(line.Text)

	switch mode {
	case style.JustifyStart:
		return helper.Right(frags, cols)
	case style.JustifyEnd:
		return helper.Left(frags, cols)
	case style.JustifyCenter, style.JustifyAround, style.JustifyEvenly:
		return helper.Center(frags, cols)
	}

	return frags
}

func TestAddGaps_SingleFragment(t *testing.T) {
	frags := text.FragmentsFromString(
		"abc",
	)

	for _, mode := range []style.Justify{
		style.JustifyStart, style.JustifyEnd, style.JustifyCenter,
		style.JustifyBetween, style.JustifyAround, style.JustifyEvenly,
	} {
		result := addGaps(10, frags, 3, mode)

		assert.Len(t, 1, result)
		assert.Equal(t, "abc", result[0].Text)
		assert.Equal(t, style.SpcKindNone, result[0].Spec.Kind())
	}
}

func TestAddGaps_IntercalatedSpaces(t *testing.T) {
	frags := text.FragmentsFromString(
		"aa", "bb", "cc",
	)

	for _, mode := range []style.Justify{style.JustifyStart, style.JustifyEnd, style.JustifyCenter} {
		result := addGaps(10, frags, 6, mode)

		assert.Len(t, 5, result)
		assert.Equal(t, "aa bb cc", fragmentTexts(result))

		assert.Equal(t, style.SpcKindNone, result[0].Spec.Kind())
		assert.Equal(t, style.SpcKindNone, result[2].Spec.Kind())
		assert.Equal(t, style.SpcKindNone, result[4].Spec.Kind())
	}
}

func TestAddGaps_Between(t *testing.T) {
	frags := text.FragmentsFromString(
		"aa", "bb", "cc",
	)

	result := addGaps(10, frags, 6, style.JustifyBetween)

	assert.Len(t, 5, result)

	assert.Equal(t, uint(2), result[1].Spec.Args()[style.KeyPaddingRightSize].Uintd(0))
	assert.Equal(t, uint(2), result[3].Spec.Args()[style.KeyPaddingRightSize].Uintd(0))
	assert.Equal(t, style.SpcKindNone, result[2].Spec.Kind())

	assert.Equal(t, "aa  bb  cc", renderFragments(result))
}

func TestAddGaps_Around(t *testing.T) {
	frags := text.FragmentsFromString(
		"aa", "bb", "cc",
	)

	result := addGaps(11, frags, 6, style.JustifyAround)

	assert.Len(t, 5, result)

	assert.Equal(t, uint(2), result[1].Spec.Args()[style.KeyPaddingRightSize].Uintd(0))
	assert.Equal(t, uint(1), result[3].Spec.Args()[style.KeyPaddingRightSize].Uintd(0))
	assert.Equal(t, style.SpcKindNone, result[2].Spec.Kind())

	assert.Equal(t, "aa  bb cc", renderFragments(result))
}

func TestAddGaps_Overflow_Start(t *testing.T) {
	frags := text.FragmentsFromString(
		"aaaa", "bbbb",
	)

	result := addGaps(5, frags, 8, style.JustifyStart)
	assert.Equal(t, "aaaa bbbb", fragmentTexts(result))
}

func TestAddGaps_DoesNotMutateOriginal(t *testing.T) {
	frags := text.FragmentsFromString(
		"aa", "bb",
	)

	_ = addGaps(10, frags, 4, style.JustifyBetween)

	for _, f := range frags {
		assert.Equal(t, style.SpcKindNone, f.Spec.Kind())
	}
}

func TestJustifyLine_Start(t *testing.T) {
	frags := text.FragmentsFromString("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyStart)

	assert.Len(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(style.SpcKindPaddingRight))

	assert.Equal(t, "aa bb cc  ", renderLine(10, style.JustifyStart, line))
}

func TestJustifyLine_End(t *testing.T) {
	frags := text.FragmentsFromString("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyEnd)

	assert.Len(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(style.SpcKindPaddingLeft))

	assert.Equal(t, "  aa bb cc", renderLine(10, style.JustifyEnd, line))
}

func TestJustifyLine_Center(t *testing.T) {
	frags := text.FragmentsFromString("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyCenter)

	assert.Len(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(style.SpcKindPaddingCenter))

	assert.Equal(t, " aa bb cc ", renderLine(10, style.JustifyCenter, line))
}

func TestJustifyLine_Between(t *testing.T) {
	frags := text.FragmentsFromString("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyBetween)

	assert.Len(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasNone(style.SpcKindPaddingLeft|style.SpcKindPaddingRight|style.SpcKindPaddingCenter))

	assert.Equal(t, "aa  bb  cc", renderLine(10, style.JustifyBetween, line))
}

func TestJustifyLine_Around(t *testing.T) {
	frags := text.FragmentsFromString("aa", "bb", "cc")
	line := justifyLine(18, frags, 6, style.JustifyAround)

	assert.Len(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(style.SpcKindPaddingCenter))

	assert.Equal(t, "   aa   bb   cc   ", renderLine(18, style.JustifyAround, line))
}

func TestJustifyLine_Evenly(t *testing.T) {
	frags := text.FragmentsFromString("aa", "bb", "cc")
	line := justifyLine(18, frags, 6, style.JustifyEvenly)

	assert.Len(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(style.SpcKindPaddingCenter))

	assert.Equal(t, "  aa    bb    cc  ", renderLine(18, style.JustifyEvenly, line))
}
