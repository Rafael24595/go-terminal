package wrapper_layout

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-terminal/engine/app/draw"
	"github.com/Rafael24595/go-terminal/engine/app/state"
	"github.com/Rafael24595/go-terminal/engine/app/viewmodel"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-terminal/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-terminal/engine/model/winsize"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"

	drawable_test "github.com/Rafael24595/go-terminal/test/engine/layout/drawable"
)

func TestTerminalApply_FixedAndPaged(t *testing.T) {
	size := winsize.Winsize{Rows: 6, Cols: 10}

	vm := viewmodel.NewViewModel()

	vm.Header.Push(
		block.BlockDrawableFromLines(
			*text.NewLine("HEADER", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	vm.Kernel.Push(
		line.LineDrawableFromLines(
			*text.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
			*text.NewLine("LINE TWO", style.SpecFromKind(style.SpcKindPaddingLeft)),
			*text.NewLine("LINE THREE IS LONG", style.SpecFromKind(style.SpcKindPaddingLeft)),
			*text.NewLine("LINE FOUR", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := text.FragmentsFromString("INPUT")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag...),
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
			width += text.FragmentMeasure(int(size.Cols), f)
		}

		assert.LessOrEqual(t, int(size.Cols), width)
	}
}

func TestTerminalApply_MultiplePages(t *testing.T) {
	size := winsize.Winsize{Rows: 4, Cols: 8}

	stt := state.NewUIState()

	vm := viewmodel.NewViewModel()

	vm.Header.Push(
		block.BlockDrawableFromLines(
			*text.NewLine("H", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	vm.Kernel.Push(
		line.LineDrawableFromLines(
			*text.NewLine("AAAAAAA", style.SpecFromKind(style.SpcKindPaddingLeft)),
			*text.NewLine("BBBBBBB", style.SpecFromKind(style.SpcKindPaddingLeft)),
			*text.NewLine("CCCCCCC", style.SpecFromKind(style.SpcKindPaddingLeft)),
			*text.NewLine("DDDDDDD", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := text.FragmentsFromString("X")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag...),
		},
	}

	vm.SetInput(viewmodel.NewInputLine(mock.ToDrawable()))

	lines0 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines0)
	assert.Equal(t, "H", lines0[0].Text[0].Text)
	assert.Equal(t, "> X", text.LineToString(&lines0[len(lines0)-1]))

	header := vm.Header.ToDrawable()
	header.Init()

	stt.Pager.TargetPage = 1
	lines1 := TerminalApply(stt, *vm, size)

	assert.Len(t, int(size.Rows), lines1)
	assert.Equal(t, "H", lines1[0].Text[0].Text)
	assert.Equal(t, "> X", text.LineToString(&lines0[len(lines0)-1]))
}

func TestDrawDynamicLines_WordWrap(t *testing.T) {
	sizeCols := 5

	lines := []text.Line{
		*text.NewLine("HELLO WORLD", style.SpecFromKind(style.SpcKindPaddingLeft)),
	}

	dw := block.BlockDrawableFromLines(lines...)

	layer := stack.NewVStackDrawable().
		Push(dw).
		ToDrawable()

	layer.Init()

	stt := state.NewUIState()

	vm := viewmodel.NewViewModel()

	dynamicSize := winsize.New(2, uint16(sizeCols))
	drawCtx := draw.NewDrawContext(stt, dynamicSize)
	drawStt := drawDynamicLines(drawCtx, vm.Pager, layer)

	assert.LessOrEqual(t, 2, len(drawStt.Buffer))

	for _, l := range drawStt.Buffer {
		width := 0
		for _, f := range l.Text {
			width += text.FragmentMeasure(sizeCols, f)
		}
		assert.LessOrEqual(t, sizeCols, width)
	}
}

func TestDrawStaticLines_DoesNotExceedRows(t *testing.T) {
	lines := []text.Line{
		*text.NewLine("golang"),
		*text.NewLine("rust"),
		*text.NewLine("ziglang"),
	}

	dw := block.BlockDrawableFromLines(lines...)

	layer := stack.NewVStackDrawable().
		Push(dw).
		ToDrawable()

	layer.Init()

	result := drawStaticLines(layer, winsize.Winsize{
		Rows: 2,
		Cols: 80,
	})

	assert.LessOrEqual(t, 2, len(result))
}

func TestDrawStaticLines_WrapThenTruncate(t *testing.T) {
	lines := []text.Line{
		*text.NewLine("golang ziglang"),
	}

	dw := block.BlockDrawableFromLines(lines...)

	layer := stack.NewVStackDrawable().
		Push(dw).
		ToDrawable()

	layer.Init()

	result := drawStaticLines(layer, winsize.Winsize{
		Rows: 3,
		Cols: 7,
	})

	assert.Len(t, 2, result)

	assert.Equal(t, "golang", result[0].Text[0].Text)
	assert.Equal(t, " ", result[0].Text[1].Text)

	assert.Equal(t, "ziglang", result[1].Text[0].Text)
}

func TestTerminalApply_InitializeLayers(t *testing.T) {
	size := winsize.Winsize{Rows: 5, Cols: 8}

	stt := state.NewUIState()

	vm := viewmodel.NewViewModel()

	vm.Header.Push(
		block.BlockDrawableFromLines(
			*text.NewLine("golang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)
	vm.Kernel.Push(
		line.LineDrawableFromLines(
			*text.NewLine("rust", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)
	vm.Footer.Push(
		block.BlockDrawableFromLines(
			*text.NewLine("Ziglang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := text.FragmentsFromString("X")
	mock := &drawable_test.MockDrawable{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag...),
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
