package wrapper_layout

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTerminalApply_FixedAndPaged(t *testing.T) {
	size := terminal.Winsize{Rows: 6, Cols: 10}

	stt := state.NewUIState()
	stt.Pager = state.PagerState{
		Page: 0,
	}

	vm := core.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.LinesEagerDrawableFromLines(
			core.NewLine("HEADER", style.SpecFromKind(style.SpcKindLeft)),
		),
	)

	vm.Lines.Shift(
		line.LinesLazyDrawableFromLines(
			core.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
			core.NewLine("LINE TWO", style.SpecFromKind(style.SpcKindLeft)),
			core.NewLine("LINE THREE IS LONG", style.SpecFromKind(style.SpcKindLeft)),
			core.NewLine("LINE FOUR", style.SpecFromKind(style.SpcKindLeft)),
		),
	)

	vm.Input = &core.InputLine{
		Prompt: ">",
		Value:  "INPUT",
	}

	state := &state.UIState{
		Pager: state.PagerState{
			Page: 0,
		},
	}

	lines := TerminalApply(state, *vm, size)

	assert.Len(t, int(size.Rows), lines)
	assert.Equal(t, "HEADER", lines[0].Text[0].Text)

	inputLine := lines[len(lines)-1]
	expectedInput := ">INPUT"

	var text strings.Builder
	for _, f := range inputLine.Text {
		text.WriteString(f.Text)
	}

	assert.Equal(t, expectedInput, text.String())

	for i := 1; i < len(lines)-1; i++ {
		width := 0
		for _, f := range lines[i].Text {
			width += utf8.RuneCountInString(f.Text)
		}

		assert.LessOrEqual(t, int(size.Cols), width)
	}
}

func TestTerminalApply_MultiplePages(t *testing.T) {
	size := terminal.Winsize{Rows: 4, Cols: 8}

	stt := state.NewUIState()
	stt.Pager = state.PagerState{
		Page: 0,
	}

	vm := core.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.LinesEagerDrawableFromLines(
			core.NewLine("H", style.SpecFromKind(style.SpcKindLeft)),
		),
	)

	vm.Lines.Shift(
		line.LinesLazyDrawableFromLines(
			core.NewLine("AAAAAAA", style.SpecFromKind(style.SpcKindLeft)),
			core.NewLine("BBBBBBB", style.SpecFromKind(style.SpcKindLeft)),
			core.NewLine("CCCCCCC", style.SpecFromKind(style.SpcKindLeft)),
			core.NewLine("DDDDDDD", style.SpecFromKind(style.SpcKindLeft)),
		),
	)

	vm.Input = &core.InputLine{
		Prompt: ">",
		Value:  "X",
	}

	lines0 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines0)
	assert.Equal(t, "H", lines0[0].Text[0].Text)
	assert.Equal(t, ">X", lines0[len(lines0)-1].Text[0].Text)

	vm.Header.Init(size)

	stt.Pager.Page = 1
	lines1 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines1)
	assert.Equal(t, "H", lines1[0].Text[0].Text)
	assert.Equal(t, ">X", lines1[len(lines1)-1].Text[0].Text)
}

func TestDrawDynamicLines_WordWrap(t *testing.T) {
	sizeCols := 5

	lines := []core.Line{
		core.NewLine("HELLO WORLD", style.SpecFromKind(style.SpcKindLeft)),
	}

	layer := core.NewLayerStack().
		Shift(line.LinesEagerDrawableFromLines(lines...))

	stt := state.NewUIState()
	stt.Pager = state.PagerState{
		Page: 0,
	}

	paged, _, _ := drawDynamicLines(stt, layer, 2, sizeCols)

	assert.LessOrEqual(t, 2, len(paged))

	for _, l := range paged {
		width := 0
		for _, f := range l.Text {
			width += utf8.RuneCountInString(f.Text)
		}
		assert.LessOrEqual(t, sizeCols, width)
	}
}

func TestDrawStaticLines_DoesNotExceedRows(t *testing.T) {
	lines := core.NewLines(
		core.LineFromString("golang"),
		core.LineFromString("rust"),
		core.LineFromString("ziglang"),
	)

	layer := core.NewLayerStack().
		Shift(line.LinesEagerDrawableFromLines(lines...))

	result := drawStaticLines(layer, 2, 80)

	assert.LessOrEqual(t, 2, len(result))
}

func TestDrawStaticLines_WrapThenTruncate(t *testing.T) {
	lines := core.NewLines(
		core.LineFromString("golang ziglang"),
	)

	layer := core.NewLayerStack().
		Shift(line.LinesEagerDrawableFromLines(lines...))

	result := drawStaticLines(layer, 3, 7)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, "golang", result[0].String())
	assert.Equal(t, " ", result[1].String())
	assert.Equal(t, "ziglang", result[2].String())
}

func TestTerminalApply_InitializeLayers(t *testing.T) {
	size := terminal.Winsize{Rows: 4, Cols: 8}

	stt := state.NewUIState()
	stt.Pager = state.PagerState{
		Page: 0,
	}

	vm := core.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.LinesEagerDrawableFromLines(
			core.NewLine("golang", style.SpecFromKind(style.SpcKindLeft)),
		),
	)
	vm.Lines.Shift(
		line.LinesLazyDrawableFromLines(
			core.NewLine("rust", style.SpecFromKind(style.SpcKindLeft)),
		),
	)
	vm.Footer.Shift(
		line.LinesEagerDrawableFromLines(
			core.NewLine("Ziglang", style.SpecFromKind(style.SpcKindLeft)),
		),
	)

	vm.Input = &core.InputLine{
		Prompt: ">",
		Value:  "X",
	}

	assert.True(t, vm.Header.HasNext())
	assert.True(t, vm.Lines.HasNext())
	assert.True(t, vm.Footer.HasNext())

	TerminalApply(stt, *vm, size)

	assert.False(t, vm.Header.HasNext())
	assert.False(t, vm.Lines.HasNext())
	assert.False(t, vm.Footer.HasNext())
}
