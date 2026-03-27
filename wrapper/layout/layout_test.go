package wrapper_layout

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/line"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stack"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestTerminalApply_FixedAndPaged(t *testing.T) {
	size := terminal.Winsize{Rows: 6, Cols: 10}

	stt := state.NewUIState()

	vm := viewmodel.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(
			text.NewLine("HEADER", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	vm.Kernel.Shift(
		line.LazyDrawableFromLines(
			text.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
			text.NewLine("LINE TWO", style.SpecFromKind(style.SpcKindPaddingLeft)),
			text.NewLine("LINE THREE IS LONG", style.SpecFromKind(style.SpcKindPaddingLeft)),
			text.NewLine("LINE FOUR", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := text.FragmentsFromString("INPUT")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			text.LineFromFragments(frag...),
		},
	}

	vm.SetInput(viewmodel.NewInputLine(mock.ToDrawable()))

	state := &state.UIState{}

	lines := TerminalApply(state, *vm, size)

	assert.Len(t, int(size.Rows), lines)
	assert.Equal(t, "HEADER", lines[0].Text[0].Text)

	inputLine := lines[len(lines)-1]
	expectedInput := "> INPUT"

	var txt strings.Builder
	for _, f := range inputLine.Text {
		txt.WriteString(f.Text)
	}

	assert.Equal(t, expectedInput, txt.String())

	for i := 1; i < len(lines)-1; i++ {
		width := 0
		for _, f := range lines[i].Text {
			width += text.FragmentMeasure(f)
		}

		assert.LessOrEqual(t, int(size.Cols), width)
	}
}

func TestTerminalApply_MultiplePages(t *testing.T) {
	size := terminal.Winsize{Rows: 4, Cols: 8}

	stt := state.NewUIState()

	vm := viewmodel.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(
			text.NewLine("H", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	vm.Kernel.Shift(
		line.LazyDrawableFromLines(
			text.NewLine("AAAAAAA", style.SpecFromKind(style.SpcKindPaddingLeft)),
			text.NewLine("BBBBBBB", style.SpecFromKind(style.SpcKindPaddingLeft)),
			text.NewLine("CCCCCCC", style.SpecFromKind(style.SpcKindPaddingLeft)),
			text.NewLine("DDDDDDD", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := text.FragmentsFromString("X")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			text.LineFromFragments(frag...),
		},
	}

	vm.SetInput(viewmodel.NewInputLine(mock.ToDrawable()))

	lines0 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines0)
	assert.Equal(t, "H", lines0[0].Text[0].Text)
	assert.Equal(t, "> X", text.LineToString(lines0[len(lines0)-1]))

	vm.Header.Init(size)

	stt.Pager.Page = 1
	lines1 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines1)
	assert.Equal(t, "H", lines1[0].Text[0].Text)
	assert.Equal(t, "> X", text.LineToString(lines0[len(lines0)-1]))
}

func TestDrawDynamicLines_WordWrap(t *testing.T) {
	sizeCols := 5

	lines := []text.Line{
		text.NewLine("HELLO WORLD", style.SpecFromKind(style.SpcKindPaddingLeft)),
	}

	dw := line.EagerDrawableFromLines(lines...)

	layer := stack.NewStackDrawable().Shift(dw)

	layer.Init(terminal.Winsize{})

	stt := state.NewUIState()

	vm := viewmodel.ViewModelFromUIState(*stt)

	paged, _, _ := drawDynamicLines(stt, *vm, layer, 2, sizeCols)

	assert.LessOrEqual(t, 2, len(paged))

	for _, l := range paged {
		width := 0
		for _, f := range l.Text {
			width += text.FragmentMeasure(f)
		}
		assert.LessOrEqual(t, sizeCols, width)
	}
}

func TestDrawStaticLines_DoesNotExceedRows(t *testing.T) {
	lines := text.NewLines(
		text.LineFromString("golang"),
		text.LineFromString("rust"),
		text.LineFromString("ziglang"),
	)

	dw := line.EagerDrawableFromLines(lines...)

	layer := stack.NewStackDrawable().Shift(dw)

	layer.Init(terminal.Winsize{})

	result := drawStaticLines(layer.ToDrawable(), 2, 80)

	assert.LessOrEqual(t, 2, len(result))
}

func TestDrawStaticLines_WrapThenTruncate(t *testing.T) {
	lines := text.NewLines(
		text.LineFromString("golang ziglang"),
	)

	dw := line.EagerDrawableFromLines(lines...)

	layer := stack.NewStackDrawable().Shift(dw)

	layer.Init(terminal.Winsize{})

	result := drawStaticLines(layer.ToDrawable(), 3, 7)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, "golang", text.LineToString(result[0]))
	assert.Equal(t, " ", text.LineToString(result[1]))
	assert.Equal(t, "ziglang", text.LineToString(result[2]))
}

func TestTerminalApply_InitializeLayers(t *testing.T) {
	size := terminal.Winsize{Rows: 5, Cols: 8}

	stt := state.NewUIState()

	vm := viewmodel.ViewModelFromUIState(*stt)

	vm.Header.Shift(
		line.EagerDrawableFromLines(
			text.NewLine("golang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)
	vm.Kernel.Shift(
		line.LazyDrawableFromLines(
			text.NewLine("rust", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)
	vm.Footer.Shift(
		line.EagerDrawableFromLines(
			text.NewLine("Ziglang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := text.FragmentsFromString("X")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			text.LineFromFragments(frag...),
		},
	}

	vm.SetInput(viewmodel.NewInputLine(mock.ToDrawable()))

	assert.True(t, vm.Header.HasNext())
	assert.True(t, vm.Kernel.HasNext())
	assert.True(t, vm.Footer.HasNext())

	TerminalApply(stt, *vm, size)

	assert.False(t, vm.Header.HasNext())
	assert.False(t, vm.Kernel.HasNext())
	assert.False(t, vm.Footer.HasNext())
}
