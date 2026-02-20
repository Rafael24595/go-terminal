package wrapper_render

import (
	"strings"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

func applyLineSpecStyles(lines []core.Line, index int, size terminal.Winsize, line string) (string, bool) {
	styl := lines[index].Spec

	baseCols := int(size.Cols)

	kind := styl.Kind()

	if kind.HasAny(style.SpcKindFill) {
		return fill(baseCols, baseCols, line), true
	}

	if kind.HasAny(style.SpcKindFillUp) {
		baseSize := baseCols
		cursor := index - 1
		if cursor >= 0 {
			baseSize = lines[cursor].Len()
		}
		return fill(baseCols, baseSize, line), true
	}

	if kind.HasAny(style.SpcKindFillDown) {
		baseSize := baseCols
		cursor := index + 1
		if cursor <= len(lines) {
			baseSize = lines[cursor].Len()
		}
		return fill(baseCols, baseSize, line), true
	}

	return line, false
}

func applySpecStyles(styl style.Spec, size terminal.Winsize, text string) string {
	baseCols := int(size.Cols)

	kind := styl.Kind()

	if kind.HasAny(style.SpcKindTrimLeft) {
		text = trimLeft(styl, text)
	}

	if kind.HasAny(style.SpcKindTrimRight) {
		text = trimRight(styl, text)
	}

	if kind.HasAny(style.SpcKindPaddingCenter) {
		text = paddingCenter(styl, baseCols, text)
	}

	if kind.HasAny(style.SpcKindPaddingLeft) {
		text = paddingLeft(styl, baseCols, text)
	}

	if kind.HasAny(style.SpcKindPaddingRight) {
		text = paddingRight(styl, baseCols, text)
	}

	if kind.HasAny(style.SpcKindRepeatLeft) {
		text = repeatLeft(styl, baseCols, text)
	}

	if kind.HasAny(style.SpcKindRepeatRight) {
		text = repeatRight(styl, baseCols, text)
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

func paddingCenter(styl style.Spec, cols int, data string) string {
	args := styl.Args()

	size := args[style.KeyPaddingCenterSize].Intd(cols)
	text := args[style.KeyPaddingCenterText].String()

	return helper.CenterCustom(data, min(cols, size), text)
}

func paddingLeft(styl style.Spec, cols int, data string) string {
	args := styl.Args()

	size := args[style.KeyPaddingLeftSize].Intd(cols)
	text := args[style.KeyPaddingLeftText].String()

	return helper.LeftCustom(data, min(cols, size), text)
}

func paddingRight(styl style.Spec, cols int, data string) string {
	args := styl.Args()

	size := args[style.KeyPaddingRightSize].Intd(cols)
	text := args[style.KeyPaddingRightText].String()

	return helper.RightCustom(data, min(cols, size), text)
}

func repeatLeft(styl style.Spec, cols int, data string) string {
	args := styl.Args()

	size := args[style.KeyRepeatLeftSize].Intd(0)
	text := args[style.KeyRepeatLeftText].String()

	if text == "" {
		text = data
		data = ""
	}

	return helper.RepeatLeftCustom(data, min(cols, size), text)
}

func repeatRight(styl style.Spec, cols int, data string) string {
	args := styl.Args()

	size := args[style.KeyRepeatRightSize].Intd(0)
	text := args[style.KeyRepeatRightText].String()

	if text == "" {
		text = data
		data = ""
	}

	return helper.RepeatRightCustom(data, min(cols, size), text)
}

func trimLeft(styl style.Spec, data string) string {
	if data == "" {
		return data
	}

	args := styl.Args()

	size := args[style.KeyTrimLeftSize].Intd(0)
	size = max(1, size)

	if size >= utf8.RuneCountInString(data) {
		return data
	}

	return data[size:]
}

func trimRight(styl style.Spec, data string) string {
	if data == "" {
		return data
	}

	args := styl.Args()

	size := args[style.KeyTrimRightSize].Intd(0)
	size = max(1, size)

	if size >= utf8.RuneCountInString(data) {
		return data
	}

	return data[:size]
}

func fill(cols, size int, data string) string {
	return helper.FillRight(data, min(cols, size))
}
