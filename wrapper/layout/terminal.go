package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm core.ViewModel, size terminal.Winsize) []core.Line {
	headerLines := make([]core.Line, 0)
	for _, header := range vm.Header {
		if header.Len() > int(size.Cols) {
			headerLines = append(headerLines, splitLineWords(int(size.Cols), header)...)
		} else {
			headerLines = append(headerLines, header)
		}
	}

	footerLines := make([]core.Line, 0)
	for _, footer := range vm.Footer {
		if footer.Len() > int(size.Cols) {
			footerLines = append(footerLines, splitLineWords(int(size.Cols), footer)...)
		} else {
			footerLines = append(footerLines, footer)
		}
	}

	inputLines := make([]core.Line, 0)
	if vm.Input != nil {
		inputLine := core.NewLine(vm.Input.Prompt+vm.Input.Value, core.ModePadding(core.Left))
		if inputLine.Len() > int(size.Cols) {
			inputLines = append(inputLines, splitLineWords(int(size.Cols), inputLine)...)
		} else {
			inputLines = append(inputLines, inputLine)
		}
	}

	bodyLines := make([]core.Line, 0)
	for _, line := range vm.Lines {
		if line.Len() > int(size.Cols) {
			bodyLines = append(bodyLines, splitLineWords(int(size.Cols), line)...)
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	rest := int(size.Rows) - (len(headerLines) + len(footerLines) + len(inputLines))
	if rest < 0 {
		return core.NewLines(
			core.LineFromString("Too low resolution"),
		)
	}

	bodyLines, page, pagination := terminalApplyBuffer(state, bodyLines, rest, int(size.Cols))

	state.Pager.Page = page
	state.Pager.Enabled = pagination

	allLines := headerLines
	allLines = append(allLines, bodyLines...)
	allLines = append(allLines, footerLines...)
	allLines = append(allLines, inputLines...)

	return allLines
}

func terminalApplyBuffer(state *state.UIState, lines []core.Line, rows, cols int) ([]core.Line, uint, bool) {
	page := uint(0)
	row := make([]core.Line, rows)

	rowCursor := 0
	lineCursor := 0
	sourceTotal := uint(0)

	for lineCursor < len(lines) {
		line := lines[lineCursor]
		lineLen := uint(line.Len())

		var fixedLines []core.Line

		if line.Len() > cols {
			fixedLines = splitLineWords(cols, line)
		} else {
			fixedLines = core.NewLines(line)
		}

		for _, fixedLine := range fixedLines {
			row[rowCursor] = fixedLine

			rowCursor += 1
			if rowCursor != rows {
				continue
			}

			isCustomFocus := state.Pager.Enabled || state.Cursor.Enabled

			isPage := state.Pager.Enabled && page == state.Pager.Page
			isCursor := state.Cursor.Enabled && sourceTotal+lineLen >= state.Cursor.Cursor

			if !isCustomFocus || isPage || isCursor {
				pagination := lineCursor != len(lines) || page != 0
				return row, page, pagination
			}

			rowCursor = 0
			row = make([]core.Line, rows)

			page++
		}

		sourceTotal += lineLen
		lineCursor++
	}

	pagination := lineCursor != len(lines) || page != 0

	return row, page, pagination
}

func splitLineWords(cols int, line core.Line) []core.Line {
	result := make([]core.Line, 0)
	current := core.LineFromPadding(line.Padding)
	width := 0

	words := core.TokenizeLineWords(line)

	for _, word := range words {
		wordlen := word.Size()

		if width+wordlen <= cols {
			current.Text = append(current.Text, word.Text...)
			width += wordlen

			continue
		}

		if wordlen <= cols {
			result = append(result, current)
			current = core.LineFromPadding(line.Padding)

			current.Text = append(current.Text, word.Text...)
			width = wordlen

			continue
		}

		newCurrent, lines, newWidth := core.SplitLongToken(word, cols, current, width)

		result = append(result, lines...)
		current = newCurrent
		width = newWidth
	}

	if len(current.Text) > 0 {
		result = append(result, current)
	}

	return result
}
