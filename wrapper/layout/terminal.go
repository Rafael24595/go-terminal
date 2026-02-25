package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm core.ViewModel, size terminal.Winsize) []core.Line {
	rows := int(size.Rows)
	cols := int(size.Cols)

	header, footer := vm.InitStaticLayers(size)

	headerLines := drawStaticLines(header, rows, cols)
	footerLines := drawStaticLines(footer, rows, cols)

	inputLines := make([]core.Line, 0)
	if input, ok := vm.InitInputLine(size); ok {
		inputLines = drawLines(*input, rows, cols)
	}

	rest := int(size.Rows) - (len(headerLines) + len(footerLines) + len(inputLines))
	if rest < 0 {
		return core.NewLines(
			core.LineFromString("Too low resolution"),
		)
	}

	remSize := terminal.NewWinsize(uint16(rest), size.Cols)
	lines := vm.InitDynamicLayers(remSize)

	bodyLines, page, remains := drawDynamicLines(state, vm, lines, rest, int(size.Cols))

	state.Pager.Page = page
	state.Pager.RestData = remains

	allLines := headerLines
	allLines = append(allLines, bodyLines...)
	allLines = append(allLines, footerLines...)
	allLines = append(allLines, inputLines...)

	return allLines
}

func drawLines(drawable core.Drawable, rows, cols int) []core.Line {
	buffer := make([]core.Line, 0)

	content := true
	for content {
		lines, status := drawable.Draw()
		content = status

		if len(lines) == 0 {
			break
		}

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

func drawDynamicLines(stt *state.UIState, vm core.ViewModel, layer *core.LayerStack, rows, cols int) ([]core.Line, uint, bool) {
	buffer := make([]core.Line, rows)
	page := uint(0)

	if rows <= 0 {
		return buffer, page, false
	}

	row := 0
	runes := uint(0)

	focus := false

	for lines := range layer.Iterator() {
		for i, line := range lines {
			lineRunes := uint(max(1, core.LineFragmentsMeasure(line)))

			fixed := fixLineSize(line, cols)
			for j, v := range fixed {
				buffer[row] = v

				if f := core.HasFocus(v); f {
					focus = f
				}

				row += 1
				if row != rows {
					continue
				}

				matches := vm.PagerMatch(*stt, state.PagerContext{
					Page:   page,
					Cursor: runes + lineRunes,
					Focus:  focus,
				})

				if matches {
					hasSpace := i < len(lines)-1 || j < len(fixed)-1
					pagination := hasSpace || layer.HasNext() || page != 0
					return buffer, page, pagination
				}

				row = 0
				buffer = make([]core.Line, rows)
				focus = false

				page++
			}
			runes += lineRunes
		}
	}

	pagination := layer.HasNext() || page != 0
	return buffer, page, pagination
}

func fixLineSize(lin core.Line, col int) []core.Line {
	if col >= core.LineFragmentsMeasure(lin) {
		return []core.Line{lin}
	}
	return line.WrapLineWords(int(col), lin)
}
