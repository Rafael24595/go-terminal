package composer

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestStandard_FixedAndPaged(t *testing.T) {
	size := winsize.Winsize{Rows: 6, Cols: 10}

	vm := viewmodel.NewViewModel()

	vm.Header.Push(
		drain.UnitFromLines(
			*text.NewLine("HEADER", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	vm.Kernel.Push(
		line.UnitFromLines(
			*text.NewLine("=", style.SpecFromKind(style.SpcKindFill)),
			*text.NewLine("LINE TWO", style.SpecFromKind(style.SpcKindPaddingLeft)),
			*text.NewLine("LINE THREE IS LONG", style.SpecFromKind(style.SpcKindPaddingLeft)),
			*text.NewLine("LINE FOUR", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
	)

	frag := text.FragmentsFromString("INPUT")
	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag...),
		},
	}

	vm.Footer.Unshift(
		inputline.Wrap(
			mock.ToUnit(),
		),
	)

	state := &state.UIState{}

	_, lines := Standard(state, *vm, size)

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
		width := winsize.Cols(0)
		for _, f := range lines[i].Text {
			width += text.FragmentMeasure(size.Cols, f)
		}

		assert.LessOrEqual(t, size.Cols, width)
	}
}

func TestStandard_InitializeLayers(t *testing.T) {
	size := winsize.Winsize{Rows: 5, Cols: 8}

	stt := state.NewUIState()

	vm := viewmodel.NewViewModel()

	vm.Header.PushLayer(
		drain.UnitFromLines(
			*text.NewLine("golang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
		layer.Fixed[winsize.Rows](1),
	)
	vm.Kernel.PushLayer(
		line.UnitFromLines(
			*text.NewLine("rust", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
		layer.Fixed[winsize.Rows](1),
	)
	vm.Footer.PushLayer(
		drain.UnitFromLines(
			*text.NewLine("Ziglang", style.SpecFromKind(style.SpcKindPaddingLeft)),
		),
		layer.Fixed[winsize.Rows](1),
	)

	frag := text.FragmentsFromString("X")
	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []text.Line{
			*text.LineFromFragments(frag...),
		},
	}

	vm.Footer.Unshift(
		inputline.Wrap(
			mock.ToUnit(),
		),
	)

	assert.True(t, vm.Header.HasNext())
	assert.True(t, vm.Kernel.HasNext())
	assert.True(t, vm.Footer.HasNext())

	Standard(stt, *vm, size)

	assert.False(t, vm.Header.HasNext())
	assert.False(t, vm.Kernel.HasNext())
	assert.False(t, vm.Footer.HasNext())
}
