package wrapper_render

import (
	"strings"

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

func applySpecStyles(styl style.Spec, size terminal.Winsize, line string) (string, bool) {
	baseCols := int(size.Cols)

	kind := styl.Kind()

	if kind.HasAny(style.SpcKindPaddingCenter) {
		return paddingCenter(styl, baseCols, line), true
	}

	if kind.HasAny(style.SpcKindPaddingLeft) {
		return paddingLeft(styl, baseCols, line), true
	}

	if kind.HasAny(style.SpcKindPaddingRight) {
		return paddingRight(styl, baseCols, line), true
	}

	if kind.HasAny(style.SpcKindRepeatLeft) {
		return repeatLeft(styl, baseCols, line), true
	}

	if kind.HasAny(style.SpcKindRepeatRight) {
		return repeatRight(styl, baseCols, line), true
	}

	return line, false
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

func fill(cols, size int, data string) string {
	return helper.Fill(data, min(cols, size))
}
