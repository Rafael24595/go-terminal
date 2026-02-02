package wrapper_layout

import (
	"testing"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

func TestSplitLineWords_Simple(t *testing.T) {
	line := core.NewLine(
		"HELLO WORLD",
		core.ModePadding(core.Left),
	)

	maxWidth := 5
	lines := splitLineWords(maxWidth, line)

	expected := []string{"HELLO", " " ,"WORLD"}

	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}

	for i, l := range lines {
		text := ""
		for _, f := range l.Text {
			text += f.Text
		}
		if text != expected[i] {
			t.Errorf("line %d expected '%s', got '%s'", i, expected[i], text)
		}
	}
}

func TestSplitLineWords_Styles(t *testing.T) {
	line := core.FragmentLine(
		core.ModePadding(core.Left),
		core.NewFragment("HELLO", core.Bold),
		core.NewFragment("WORLD"),
	)

	maxWidth := 7
	lines := splitLineWords(maxWidth, line)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	if lines[0].Text[0].Text != "HELLO" || lines[0].Text[0].Styles[0] != core.Bold {
		t.Errorf("first fragment incorrect: %+v", lines[0].Text[0])
	}
	if lines[1].Text[0].Text != "WORLD" {
		t.Errorf("second fragment incorrect: %+v", lines[0].Text[1])
	}
}

func TestSplitLineWords_LongWord(t *testing.T) {
	text := "HELLO WORLD FROM GOLANG"

	line := core.NewLine(
		text,
		core.ModePadding(core.Left),
	)

	maxWidth := 10
	lines := splitLineWords(maxWidth, line)

	for i, l := range lines {
		text := ""
		for _, f := range l.Text {
			text += f.Text
		}
		if utf8.RuneCountInString(text) > maxWidth {
			t.Errorf("line %d too long: %s", i, text)
		}
	}

	totalRunes := 0
	for _, l := range lines {
		for _, f := range l.Text {
			totalRunes += utf8.RuneCountInString(f.Text)
		}
	}
	if totalRunes != utf8.RuneCountInString(text) {
		t.Errorf("total runes mismatch")
	}
}

func TestSplitLineWords_MultipleFragments(t *testing.T) {
	line := core.FragmentLine(
		core.ModePadding(core.Left),
		core.NewFragment("HELLO", core.Bold),
		core.NewFragment("WORLD", core.Bold),
		core.NewFragment("GO"),
	)

	maxWidth := 8
	lines := splitLineWords(maxWidth, line)

	for _, l := range lines {
		width := 0
		for _, f := range l.Text {
			width += utf8.RuneCountInString(f.Text)
		}
		if width > maxWidth {
			t.Errorf("line exceeds maxWidth: %v", l)
		}
	}
}

func TestTerminalApply_FixedAndPaged(t *testing.T) {
	size := terminal.Winsize{Rows: 6, Cols: 10}

	vm := core.ViewModel{
		Header: []core.Line{
			core.NewLine("HEADER", core.ModePadding(core.Left)),
		},
		Lines: []core.Line{
			core.NewLine("=", core.ModePadding(core.Fill)),
			core.NewLine("LINE TWO", core.ModePadding(core.Left)),
			core.NewLine("LINE THREE IS LONG", core.ModePadding(core.Left)),
			core.NewLine("LINE FOUR", core.ModePadding(core.Left)),
		},
		Input: &core.InputLine{
			Prompt: ">",
			Value:  "INPUT",
		},
	}

	state := &state.UIState{
		Pager: state.PagerState{
			Page: 0,
		},
	}

	lines := TerminalApply(state, vm, size)

	if len(lines) != int(size.Rows) {
		t.Errorf("expected %d lines, got %d", size.Rows, len(lines))
	}

	if lines[0].Text[0].Text != "HEADER" {
		t.Errorf("expected header line, got %v", lines[0].Text)
	}

	inputLine := lines[len(lines)-1]
	expectedInput := ">INPUT"
	text := ""
	for _, f := range inputLine.Text {
		text += f.Text
	}
	if text != expectedInput {
		t.Errorf("expected input line '%s', got '%s'", expectedInput, text)
	}

	for i := 1; i < len(lines)-1; i++ {
		width := 0
		for _, f := range lines[i].Text {
			width += utf8.RuneCountInString(f.Text)
		}
		if width > int(size.Cols) {
			t.Errorf("line %d exceeds width %d: %+v", i, size.Cols, lines[i])
		}
	}
}

func TestTerminalApply_MultiplePages(t *testing.T) {
	size := terminal.Winsize{Rows: 4, Cols: 8}

	vm := core.ViewModel{
		Header: []core.Line{
			core.NewLine("H", core.ModePadding(core.Left)),
		},
		Lines: []core.Line{
			core.NewLine("AAAAAAA", core.ModePadding(core.Left)),
			core.NewLine("BBBBBBB", core.ModePadding(core.Left)),
			core.NewLine("CCCCCCC", core.ModePadding(core.Left)),
			core.NewLine("DDDDDDD", core.ModePadding(core.Left)),
		},
		Input: &core.InputLine{
			Prompt: ">",
			Value:  "X",
		},
	}

	state := &state.UIState{
		Pager: state.PagerState{
			Page: 0,
		},
	}

	lines0 := TerminalApply(state, vm, size)
	if len(lines0) != int(size.Rows) {
		t.Errorf("page 0: expected %d lines, got %d", size.Rows, len(lines0))
	}
	if lines0[0].Text[0].Text != "H" {
		t.Errorf("page 0: header mismatch")
	}
	if lines0[len(lines0)-1].Text[0].Text != ">X" {
		t.Errorf("page 0: input mismatch")
	}

	state.Pager.Page = 1
	lines1 := TerminalApply(state, vm, size)
	if len(lines1) != int(size.Rows) {
		t.Errorf("page 1: expected %d lines, got %d", size.Rows, len(lines1))
	}
	if lines1[0].Text[0].Text != "H" {
		t.Errorf("page 1: header mismatch")
	}
	if lines1[len(lines1)-1].Text[0].Text != ">X" {
		t.Errorf("page 1: input mismatch")
	}
}

func TestTerminalApplyBuffer_WordWrap(t *testing.T) {
	sizeCols := 5
	lines := []core.Line{
		core.NewLine("HELLO WORLD", core.ModePadding(core.Left)),
	}

	state := &state.UIState{
		Pager: state.PagerState{
			Page: 0,
		},
	}

	paged, _, _ := terminalApplyBuffer(state, lines, 2, sizeCols)

	if len(paged) > 2 {
		t.Errorf("expected max 2 lines for buffer, got %d", len(paged))
	}

	for _, l := range paged {
		width := 0
		for _, f := range l.Text {
			width += utf8.RuneCountInString(f.Text)
		}
		if width > sizeCols {
			t.Errorf("line exceeds width: %+v", l)
		}
	}
}
