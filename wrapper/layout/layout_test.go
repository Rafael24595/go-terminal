package wrapper_layout

import (
	"strings"
	"testing"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/core"
	"github.com/Rafael24595/go-terminal/engine/core/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/core/style"
	"github.com/Rafael24595/go-terminal/engine/terminal"
	drawable_test "github.com/Rafael24595/go-terminal/test/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestTerminalApply_FixedAndPaged(t *testing.T) {
	size := terminal.Winsize{Rows: 6, Cols: 10}

	stt := state.NewUIState()

	vm := core.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(
			core.NewLine("HEADER", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	vm.Lines.Shift(
		line.LazyDrawableFromLines(
			core.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
			core.NewLine("LINE TWO", style.SpecFromKind(style.SpcKindPaddingLeft)),
			core.NewLine("LINE THREE IS LONG", style.SpecFromKind(style.SpcKindPaddingLeft)),
			core.NewLine("LINE FOUR", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := core.FragmentsFromString("INPUT")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []core.Line{
			core.LineFromFragments(frag...),
		},
	}

	vm.SetInput(core.NewInputLine(mock.ToDrawable()))

	state := &state.UIState{}

	lines := TerminalApply(state, *vm, size)

	assert.Len(t, int(size.Rows), lines)
	assert.Equal(t, "HEADER", lines[0].Text[0].Text)

	inputLine := lines[len(lines)-1]
	expectedInput := "> INPUT"

	var text strings.Builder
	for _, f := range inputLine.Text {
		text.WriteString(f.Text)
	}

	assert.Equal(t, expectedInput, text.String())

	for i := 1; i < len(lines)-1; i++ {
		width := 0
		for _, f := range lines[i].Text {
			width += core.FragmentMeasure(f)
		}

		assert.LessOrEqual(t, int(size.Cols), width)
	}
}

func TestTerminalApply_MultiplePages(t *testing.T) {
	size := terminal.Winsize{Rows: 4, Cols: 8}

	stt := state.NewUIState()

	vm := core.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(
			core.NewLine("H", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	vm.Lines.Shift(
		line.LazyDrawableFromLines(
			core.NewLine("AAAAAAA", style.SpecFromKind(style.SpcKindPaddingLeft)),
			core.NewLine("BBBBBBB", style.SpecFromKind(style.SpcKindPaddingLeft)),
			core.NewLine("CCCCCCC", style.SpecFromKind(style.SpcKindPaddingLeft)),
			core.NewLine("DDDDDDD", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := core.FragmentsFromString("X")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []core.Line{
			core.LineFromFragments(frag...),
		},
	}

	vm.SetInput(core.NewInputLine(mock.ToDrawable()))

	lines0 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines0)
	assert.Equal(t, "H", lines0[0].Text[0].Text)
	assert.Equal(t, "> X", core.LineToString(lines0[len(lines0)-1]))

	vm.Header.Init(size)

	stt.Pager.Page = 1
	lines1 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines1)
	assert.Equal(t, "H", lines1[0].Text[0].Text)
	assert.Equal(t, "> X", core.LineToString(lines0[len(lines0)-1]))
}

func TestDrawDynamicLines_WordWrap(t *testing.T) {
	sizeCols := 5

	lines := []core.Line{
		core.NewLine("HELLO WORLD", style.SpecFromKind(style.SpcKindPaddingLeft)),
	}

	dw := line.EagerDrawableFromLines(lines...)

	dw.Init(terminal.Winsize{})

	layer := core.NewLayerStack().Shift(dw)

	stt := state.NewUIState()

	vm := core.ViewModelFromUIState(*stt)

	paged, _, _ := drawDynamicLines(stt, *vm, layer, 2, sizeCols)

	assert.LessOrEqual(t, 2, len(paged))

	for _, l := range paged {
		width := 0
		for _, f := range l.Text {
			width += core.FragmentMeasure(f)
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

	dw := line.EagerDrawableFromLines(lines...)

	dw.Init(terminal.Winsize{})

	layer := core.NewLayerStack().Shift(dw)

	result := drawStaticLines(layer, 2, 80)

	assert.LessOrEqual(t, 2, len(result))
}

func TestDrawStaticLines_WrapThenTruncate(t *testing.T) {
	lines := core.NewLines(
		core.LineFromString("golang ziglang"),
	)

	dw := line.EagerDrawableFromLines(lines...)

	dw.Init(terminal.Winsize{})

	layer := core.NewLayerStack().Shift(dw)

	result := drawStaticLines(layer, 3, 7)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, "golang", core.LineToString(result[0]))
	assert.Equal(t, " ", core.LineToString(result[1]))
	assert.Equal(t, "ziglang", core.LineToString(result[2]))
}

func TestTerminalApply_InitializeLayers(t *testing.T) {
	size := terminal.Winsize{Rows: 5, Cols: 8}

	stt := state.NewUIState()

	vm := core.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(
			core.NewLine("golang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)
	vm.Lines.Shift(
		line.LazyDrawableFromLines(
			core.NewLine("rust", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)
	vm.Footer.Shift(
		line.EagerDrawableFromLines(
			core.NewLine("Ziglang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := core.FragmentsFromString("X")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []core.Line{
			core.LineFromFragments(frag...),
		},
	}

	vm.SetInput(core.NewInputLine(mock.ToDrawable()))

	assert.True(t, vm.Header.HasNext())
	assert.True(t, vm.Lines.HasNext())
	assert.True(t, vm.Footer.HasNext())

	TerminalApply(stt, *vm, size)

	assert.False(t, vm.Header.HasNext())
	assert.False(t, vm.Lines.HasNext())
	assert.False(t, vm.Footer.HasNext())
}
