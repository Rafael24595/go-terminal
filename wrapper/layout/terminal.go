package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm core.ViewModel, size terminal.Winsize) []text.Line {
	rows := int(size.Rows)
	cols := int(size.Cols)

	header, footer := vm.InitStaticLayers(size)

	headerLines := drawStaticLines(header.ToDrawable(), rows, cols)
	footerLines := drawStaticLines(footer.ToDrawable(), rows, cols)

	inputLines := make([]text.Line, 0)
	if input, ok := vm.InitInputLine(size); ok {
		inputLines = drawStaticLines(input, rows, cols)
	}

	helperLines := make([]text.Line, 0)
	if helper, ok := vm.InitHelper(size); ok {
		helperLines = drawStaticLines(helper, rows, cols)
	}

	static := len(headerLines) + len(footerLines) + len(inputLines) + len(helperLines)
	rest := int(size.Rows) - static
	if rest < 0 {
		return text.NewLines(
			text.LineFromString("Too low resolution"),
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
	allLines = append(allLines, helperLines...)

	return allLines
}

func drawStaticLines(drawable drawable.Drawable, rows, cols int) []text.Line {
	buffer := make([]text.Line, 0)

	content := true
	for content {
		lines, status := drawable.Draw()
		content = status

		if len(lines) == 0 {
			break
		}

		for _, lin := range lines {
			buffer = append(buffer,
				line.WrapLineWords(cols, lin)...,
			)

			if len(buffer) >= rows {
				break
			}
		}
	}

	return buffer
}

func drawDynamicLines(stt *state.UIState, vm core.ViewModel, layer *stack.StackDrawable, rows, cols int) ([]text.Line, uint, bool) {
	buffer := make([]text.Line, rows)
	page := uint(0)

	if rows <= 0 {
		return buffer, page, false
	}

	row := 0
	runes := uint(0)

	focus := false

	for lines := range layer.Iterator() {
		for i, lin := range lines {
			lineRunes := uint(max(1, text.LineFragmentsMeasure(lin)))

			fixed := line.WrapLineWords(cols, lin)
			for j, v := range fixed {
				buffer[row] = v

				if f := text.HasFocus(v); f {
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
				buffer = make([]text.Line, rows)
				focus = false

				page++
			}
			runes += lineRunes
		}
	}

	pagination := layer.HasNext() || page != 0
	return buffer, page, pagination
}
