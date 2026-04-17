package wrapper_render

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/dict"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

var specStylesTable = dict.NewInmutableLinkedMap(
	dict.P(style.SpcKindFill, func(spec style.Spec, cols int, text string, logicalSize int) (string, bool) {
		return fill(spec, cols, text, logicalSize), true
	}),
	dict.P(style.SpcKindTrimLeft, func(spec style.Spec, _ int, text string, logicalSize int) (string, bool) {
		return trimLeft(spec, text, logicalSize), false
	}),
	dict.P(style.SpcKindTrimRight, func(spec style.Spec, _ int, text string, logicalSize int) (string, bool) {
		return trimRight(spec, text, logicalSize), false
	}),
	dict.P(style.SpcKindPaddingCenter, func(spec style.Spec, cols int, text string, logicalSize int) (string, bool) {
		return paddingCenter(spec, cols, text, logicalSize), false
	}),
	dict.P(style.SpcKindPaddingLeft, func(spec style.Spec, cols int, text string, logicalSize int) (string, bool) {
		return paddingLeft(spec, cols, text, logicalSize), false
	}),
	dict.P(style.SpcKindPaddingRight, func(spec style.Spec, cols int, text string, logicalSize int) (string, bool) {
		return paddingRight(spec, cols, text, logicalSize), false
	}),
	dict.P(style.SpcKindRepeatLeft, func(spec style.Spec, cols int, text string, logicalSize int) (string, bool) {
		return repeatLeft(spec, cols, text, logicalSize), false
	}),
	dict.P(style.SpcKindRepeatRight, func(spec style.Spec, cols int, text string, logicalSize int) (string, bool) {
		return repeatRight(spec, cols, text, logicalSize), false
	}),
)

var specAtomTable = dict.NewInmutableLinkedMap(
	dict.P(style.AtmLower, func(text string) string {
		return strings.ToLower(text)
	}),
	dict.P(style.AtmUpper, func(text string) string {
		return strings.ToUpper(text)
	}),
	dict.P(style.AtmBold, func(text string) string {
		return Bold + text + NoBold
	}),
	dict.P(style.AtmSelect, func(text string) string {
		return Reverse + text + NoReverse
	}),
)

func applySpecStyles(spec style.Spec, size terminal.Winsize, text string, logicalSize int) string {
	var exit bool
	var cols int = int(size.Cols)

	kind := spec.Kind()
	for k, p := range specStylesTable.All() {
		if !kind.HasAny(k) {
			continue
		}

		text, exit = p(spec, cols, text, logicalSize)
		if exit {
			return text
		}

		logicalSize = style.SpecMeasureOf(k, spec, style.LayoutContext{
			Text: logicalSize,
			Cols: cols,
		})
	}

	return text
}

func applyAtomStyles(text string, styles ...style.Atom) string {
	merged := style.MergeAtom(styles...)

	for k, p := range specAtomTable.All() {
		if merged.HasAny(k) {
			text = p(text)
		}
	}

	return text
}

func fill(styl style.Spec, cols int, data string, logicalSize int) string {
	opts := helper.LogicalSizeOpts{
		LogicalSize: logicalSize,
	}

	args := styl.Args()

	size := args[style.KeyFillSize].Intd(cols)
	size = min(cols, size)

	return helper.FillRightWithOpts(data, min(cols, size), opts)
}

func trimLeft(styl style.Spec, data string, logicalSize int) string {
	if data == "" {
		return data
	}

	args := styl.Args()

	size := args[style.KeyTrimLeftSize].Intd(0)
	size = max(1, size)

	elip := args[style.KeyTrimEllipsisText].Stringf()

	opts := helper.TextTrimOpts{
		LogicalSize:  logicalSize,
		EllipsisText: elip,
	}

	return helper.TrimLeft(data, size, opts)
}

func trimRight(styl style.Spec, data string, logicalSize int) string {
	if data == "" {
		return data
	}

	args := styl.Args()

	size := args[style.KeyTrimRightSize].Intd(0)
	size = max(1, size)

	elip := args[style.KeyTrimEllipsisText].Stringf()

	opts := helper.TextTrimOpts{
		LogicalSize:  logicalSize,
		EllipsisText: elip,
	}

	return helper.TrimRight(data, size, opts)
}

// TODO: Explore the risks of using cols as default
func paddingCenter(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyPaddingCenterSize].Intd(cols)
	text := args[style.KeyPaddingCenterText].Stringf()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Runes:       text,
	}

	return helper.CenterWithOpts(data, min(cols, size), opts)
}

func paddingLeft(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyPaddingLeftSize].Intd(0)
	text := args[style.KeyPaddingLeftText].Stringf()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Runes:       text,
	}

	return helper.LeftWithOpts(data, min(cols, size), opts)
}

func paddingRight(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyPaddingRightSize].Intd(0)
	text := args[style.KeyPaddingRightText].Stringf()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Runes:       text,
	}

	return helper.RightWithOpts(data, min(cols, size), opts)
}

func repeatLeft(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyRepeatLeftSize].Intd(0)
	text := args[style.KeyRepeatLeftText].Stringf()

	if text == "" {
		text = data
		data = ""
	}

	opts := helper.LogicalSizeOpts{
		LogicalSize: logicalSize,
	}

	return helper.RepeatLeftWithOpts(data, text, min(cols, size), opts)
}

func repeatRight(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyRepeatRightSize].Intd(0)
	text := args[style.KeyRepeatRightText].Stringf()

	if text == "" {
		text = data
		data = ""
	}

	opts := helper.LogicalSizeOpts{
		LogicalSize: logicalSize,
	}

	return helper.RepeatRightWithOpts(data, text, min(cols, size), opts)
}
