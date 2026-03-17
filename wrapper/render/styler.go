package wrapper_render

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

func applySpecStyles(styl style.Spec, size terminal.Winsize, text string, logicalSize int) string {
	baseCols := int(size.Cols)

	kind := styl.Kind()

	if kind.HasAny(style.SpcKindFill) {
		return fill(styl, baseCols, text, logicalSize)
	}

	if kind.HasAny(style.SpcKindTrimLeft) {
		text = trimLeft(styl, text, logicalSize)
	}

	if kind.HasAny(style.SpcKindTrimRight) {
		text = trimRight(styl, text, logicalSize)
	}

	if kind.HasAny(style.SpcKindPaddingCenter) {
		text = paddingCenter(styl, baseCols, text, logicalSize)
	}

	if kind.HasAny(style.SpcKindPaddingLeft) {
		text = paddingLeft(styl, baseCols, text, logicalSize)
	}

	if kind.HasAny(style.SpcKindPaddingRight) {
		text = paddingRight(styl, baseCols, text, logicalSize)
	}

	if kind.HasAny(style.SpcKindRepeatLeft) {
		text = repeatLeft(styl, baseCols, text, logicalSize)
	}

	if kind.HasAny(style.SpcKindRepeatRight) {
		text = repeatRight(styl, baseCols, text, logicalSize)
	}

	return text
}

func applyAtomStyles(text string, styles ...style.Atom) string {
	merged := style.MergeAtom(styles...)

	if merged.HasAny(style.AtmLower) {
		text = strings.ToLower(text)
	}

	if merged.HasAny(style.AtmUpper) {
		text = strings.ToUpper(text)
	}

	if merged.HasAny(style.AtmBold) {
		text = Bold + text + NoBold
	}

	if merged.HasAny(style.AtmSelect) {
		text = Reverse + text + NoReverse
	}

	return text
}

func paddingCenter(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyPaddingCenterSize].Intd(cols)
	text := args[style.KeyPaddingCenterText].String()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Runes:       text,
	}

	return helper.CenterWithOpts(data, min(cols, size), opts)
}

func paddingLeft(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyPaddingLeftSize].Intd(cols)
	text := args[style.KeyPaddingLeftText].String()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Runes:       text,
	}

	return helper.LeftWithOpts(data, min(cols, size), opts)
}

func paddingRight(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyPaddingRightSize].Intd(cols)
	text := args[style.KeyPaddingRightText].String()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Runes:       text,
	}

	return helper.RightWithOpts(data, min(cols, size), opts)
}

func repeatLeft(styl style.Spec, cols int, data string, logicalSize int) string {
	args := styl.Args()

	size := args[style.KeyRepeatLeftSize].Intd(0)
	text := args[style.KeyRepeatLeftText].String()

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
	text := args[style.KeyRepeatRightText].String()

	if text == "" {
		text = data
		data = ""
	}

	opts := helper.LogicalSizeOpts{
		LogicalSize: logicalSize,
	}

	return helper.RepeatRightWithOpts(data, text, min(cols, size), opts)
}

func trimLeft(styl style.Spec, data string, logicalSize int) string {
	if data == "" {
		return data
	}

	args := styl.Args()

	size := args[style.KeyTrimLeftSize].Intd(0)
	size = max(1, size)

	elip := args[style.KeyTrimEllipsisText].String()

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

	elip := args[style.KeyTrimEllipsisText].String()

	opts := helper.TextTrimOpts{
		LogicalSize:  logicalSize,
		EllipsisText: elip,
	}

	return helper.TrimRight(data, size, opts)
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
