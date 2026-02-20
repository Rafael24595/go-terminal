package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm core.ViewModel, size terminal.Winsize) []core.Line {
	rows := int(size.Rows)
	cols := int(size.Cols)

	header, lines, footer := vm.InitLayers(size)

	headerLines := drawStaticLines(header, rows, cols)
	footerLines := drawStaticLines(footer, rows, cols)

	inputLines := make([]core.Line, 0)
	if vm.Input != nil {
		inputLine := core.NewLine(vm.Input.Prompt+vm.Input.Value, style.SpecFromKind(style.SpcKindPaddingLeft))
		if inputLine.Len() > int(size.Cols) {
			inputLines = append(inputLines, line.WrapLineWords(int(size.Cols), inputLine)...)
		} else {
			inputLines = append(inputLines, inputLine)
		}
	}

	rest := int(size.Rows) - (len(headerLines) + len(footerLines) + len(inputLines))
	if rest < 0 {
		return core.NewLines(
			core.LineFromString("Too low resolution"),
		)
	}

	bodyLines, page, pagination := drawDynamicLines(state, lines, rest, int(size.Cols))

	state.Pager.Page = page
	state.Pager.Enabled = pagination

	allLines := headerLines
	allLines = append(allLines, bodyLines...)
	allLines = append(allLines, footerLines...)
	allLines = append(allLines, inputLines...)

	return allLines
}

func drawStaticLines(layer *core.LayerStack, rows, cols int) []core.Line {
	buffer := make([]core.Line, 0)

	for lines := range layer.Iterator() {
		for _, line := range lines {
			buffer = append(buffer,
				fixLineSize(line, cols)...,
			)

			if len(buffer) >= rows {
				break
			}
		}
	}

	return buffer
}

func drawDynamicLines(state *state.UIState, layer *core.LayerStack, rows, cols int) ([]core.Line, uint, bool) {
	buffer := make([]core.Line, rows)
	page := uint(0)

	if rows <= 0 {
		return buffer, page, false
	}

	row := 0
	runes := uint(0)

	for lines := range layer.Iterator() {
		for i, line := range lines {
			lineRunes := uint(max(1, line.Len()))

			fixed := fixLineSize(line, cols)
			for j, v := range fixed {
				buffer[row] = v

				row += 1
				if row != rows {
					continue
				}

				isCustomFocus := state.Pager.Enabled || state.Cursor.Enabled

				isPage := state.Pager.Enabled && page == state.Pager.Page
				isCursor := state.Cursor.Enabled && runes+lineRunes >= state.Cursor.Cursor

				if !isCustomFocus || isPage || isCursor {
					hasSpace := i < len(lines)-1 || j < len(fixed)-1
					pagination := hasSpace || layer.HasNext() || page != 0
					return buffer, page, pagination
				}

				row = 0
				buffer = make([]core.Line, rows)

				page++
			}
			runes += lineRunes
		}
	}

	pagination := layer.HasNext() || page != 0

	return buffer, page, pagination
}

func fixLineSize(lin core.Line, col int) []core.Line {
	if col >= lin.Len() {
		return []core.Line{lin}
	}
	return line.WrapLineWords(int(col), lin)
}
