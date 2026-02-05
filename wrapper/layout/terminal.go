package wrapper_layout

import (
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/assert"
	"github.com/Rafael24595/go-terminal/engine/helper"
	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

const separator = " | "

type IndexMeta struct {
	sufix      string
	prefixBody string
	digits     uint16
	totalWidth uint32
}

func (i IndexMeta) header(index int) string {
	return helper.Right(index, int(i.digits)) + i.sufix
}

func (i IndexMeta) body() string {
	return i.prefixBody + i.sufix
}

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm core.ViewModel, size terminal.Winsize) []core.Line {
	headerLines := make([]core.Line, 0)
	for _, header := range vm.Header {
		if header.Len() > int(size.Cols) {
			headerLines = append(headerLines, wrapLineWords(int(size.Cols), header)...)
		} else {
			headerLines = append(headerLines, header)
		}
	}

	footerLines := make([]core.Line, 0)
	for _, footer := range vm.Footer {
		if footer.Len() > int(size.Cols) {
			footerLines = append(footerLines, wrapLineWords(int(size.Cols), footer)...)
		} else {
			footerLines = append(footerLines, footer)
		}
	}

	inputLines := make([]core.Line, 0)
	if vm.Input != nil {
		inputLine := core.NewLine(vm.Input.Prompt+vm.Input.Value, core.ModePadding(core.Left))
		if inputLine.Len() > int(size.Cols) {
			inputLines = append(inputLines, wrapLineWords(int(size.Cols), inputLine)...)
		} else {
			inputLines = append(inputLines, inputLine)
		}
	}

	bodyLines := make([]core.Line, 0)
	indexMeta := computeIndexMeta(vm.Lines)
	for _, line := range vm.Lines {
		bodyLines = append(bodyLines, indexLines(int(size.Cols), line, indexMeta)...)
	}

	rest := int(size.Rows) - (len(headerLines) + len(footerLines) + len(inputLines))
	if rest < 0 {
		return core.NewLines(
			core.LineFromString("Too low resolution"),
		)
	}

	bodyLines, page, pagination := terminalApplyBuffer(state, bodyLines, rest, int(size.Cols), indexMeta)

	state.Pager.Page = page
	state.Pager.Enabled = pagination

	allLines := headerLines
	allLines = append(allLines, bodyLines...)
	allLines = append(allLines, footerLines...)
	allLines = append(allLines, inputLines...)

	return allLines
}

func terminalApplyBuffer(state *state.UIState, lines []core.Line, rows, cols int, meta *IndexMeta) ([]core.Line, uint, bool) {
	page := uint(0)
	row := make([]core.Line, rows)

	rowCursor := 0
	lineCursor := 0
	sourceTotal := uint(0)

	indexFix := 0
	if meta != nil {
		indexFix = int(meta.totalWidth)
	}

	for lineCursor < len(lines) {
		line := lines[lineCursor]
		lineLen := uint (max(1, line.Len() - indexFix))

		var fixedLines []core.Line

		if line.Len() > cols {
			fixedLines = wrapLineWords(cols, line)
			assert.Unreachablef("the lines at this point should be less than cols")
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

func indexLines(cols int, line core.Line, meta *IndexMeta) []core.Line {
	isGreaterWithoutIndex := line.Len() > int(cols)
	isGreaterWithIndex := meta != nil && line.Len()+int(meta.totalWidth) > cols

	if isGreaterWithoutIndex || isGreaterWithIndex {
		return wrapLineWordsWithIndex(int(cols), line, meta)
	}

	fragments := core.FragmentsFromString()
	if meta != nil {
		fragments = append(fragments, core.NewFragment(meta.header(int(line.Order))))
	}

	newLine := core.LineFromFragments(
		append(fragments, line.Text...)...,
	)

	return core.FixedLinesFromLines(line.Padding, newLine)
}

func wrapLineWords(cols int, line core.Line) []core.Line {
	return wrapLineWordsWithIndex(cols, line, nil)
}

func wrapLineWordsWithIndex(cols int, line core.Line, meta *IndexMeta) []core.Line {
	result := make([]core.Line, 0)
	current := core.LineFromPadding(line.Padding)
	width := 0

	words := core.TokenizeLineWords(line)

	if meta != nil {
		fragments := core.FragmentsFromString(meta.header(int(line.Order)))
		current.Text = append(current.Text, fragments...)
		cols -= int(meta.totalWidth)
	}

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

			if meta != nil {
				fragments := core.FragmentsFromString(meta.body())
				current.Text = append(current.Text, fragments...)
			}

			current.Text = append(current.Text, word.Text...)
			width = wordlen

			continue
		}

		newCurrent, lines, newWidth := wrapLongTokenWithIndex(word, cols, current, width, meta)

		result = append(result, lines...)
		current = newCurrent
		width = newWidth
	}

	if len(current.Text) > 0 {
		result = append(result, current)
	}

	return result
}

func wrapLongTokenWithIndex(
	word core.WordToken,
	cols int,
	current core.Line,
	width int,
	meta *IndexMeta,
) (core.Line, []core.Line, int) {
	current, lines, width := core.SplitLongToken(word, cols, current, width)
	if meta == nil || len(lines) == 0 {
		return current, lines, width
	}

	index := core.FragmentsFromString(meta.body())

	current.Text = append(index, current.Text...)

	for i := 1; i < len(lines); i++ {
		lines[i].Text = append(index, lines[i].Text...)
	}

	return current, lines, width
}

func computeIndexMeta(lines []core.Line) *IndexMeta {
	size := uint32(0)

	for _, line := range lines {
		if line.Order == 0 {
			continue
		}
		size = max(size, math.Digits(line.Order))
	}

	if size == 0 {
		return nil
	}

	return &IndexMeta{
		sufix:      separator,
		prefixBody: helper.Fill(" ", int(size)),
		digits:     uint16(size),
		totalWidth: size + uint32(len(separator)),
	}
}
