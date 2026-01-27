package wrapper_layout

import (
	"strings"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

func TerminalApply(state *state.UIState, vm core.ViewModel, size terminal.Winsize) []core.Line {
	headerLines := make([]core.Line, 0)
	for _, header := range vm.Headers {
		if header.Len() > int(size.Cols) {
			headerLines = append(headerLines, splitLineWords(int(size.Cols), header)...)
		} else {
			headerLines = append(headerLines, header)
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

	rest := int(size.Rows) - (len(headerLines) + len(inputLines))

	bodyLines = terminalApplyBuffer(state, bodyLines, rest, int(size.Cols))

	allLines := headerLines
	allLines = append(allLines, bodyLines...)
	allLines = append(allLines, inputLines...)

	return allLines
}

func terminalApplyBuffer(state *state.UIState, lines []core.Line, rows, cols int) []core.Line {
	page := 0
	cursor := []core.Line{}

	i := 0
	for i < len(lines) {
		line := lines[i]
		var physicalLines []core.Line

		if line.Len() > cols {
			physicalLines = splitLineWords(cols, line)
		} else {
			physicalLines = []core.Line{line}
		}

		for _, pl := range physicalLines {
			cursor = append(cursor, pl)

			if len(cursor) != rows {
				continue
			}

			if page == state.Page {
				return cursor
			}

			page++
			cursor = []core.Line{}
		}

		i++
	}

	if page == state.Page {
		return cursor
	}

	return []core.Line{}
}

func splitLineWords(cols int, line core.Line) []core.Line {
	result := make([]core.Line, 0)
	current := core.LineFromPadding(line.Padding)
	width := 0

	for _, frag := range line.Text {
		words := strings.Fields(frag.Text)
		for wi, word := range words {
			space := 0
			if width > 0 && wi > 0 {
				space = 1
			}

			wordlen := utf8.RuneCountInString(word)

			if width+space+wordlen <= cols {
				current, width = appendWordToLine(current, word, frag, space, width)
				continue
			}

			if wordlen <= cols {
				result = append(result, current)
				current = core.LineFromPadding(line.Padding)
				current, width = appendWordToLine(current, word, frag, 0, 0)
				continue
			}

			lines, newCurrent, newWidth := splitLongWord(word, frag, cols, current, width)
			result = append(result, lines...)
			current = newCurrent
			width = newWidth
		}
	}

	if len(current.Text) > 0 {
		result = append(result, current)
	}

	return result
}

func appendWordToLine(line core.Line, word string, frag core.Fragment, space int, width int) (core.Line, int) {
	if space > 0 {
		fragment := core.NewFragment(" ")
		line.Text = append(line.Text, fragment)
		width += 1
	}

	fragment := core.NewFragment(word, frag.Styles...)
	line.Text = append(line.Text, fragment)

	return line, width + utf8.RuneCountInString(word)
}

func splitLongWord(word string, frag core.Fragment, cols int, current core.Line, width int) ([]core.Line, core.Line, int) {
	result := make([]core.Line, 0)
	runes := []rune(word)
	start := 0

	for start < len(runes) {
		remaining := cols - width
		if remaining <= 0 {
			result = append(result, current)
			current = core.LineFromPadding(current.Padding)
			width = 0
			remaining = cols
		}

		end := start + remaining
		if end > len(runes) {
			end = len(runes)
		}

		word := string(runes[start:end])
		fragment := core.NewFragment(word, frag.Styles...)
		current.Text = append(current.Text, fragment)

		width += utf8.RuneCountInString(word)
		start = end
	}

	return result, current, width
}
