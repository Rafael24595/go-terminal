package composer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func Standard(
	uiState *state.UIState,
	vm viewmodel.ViewModel,
	size winsize.Winsize,
) (*state.UIState, []text.Line) {
	header := vm.Header.ToUnit()
	header.Drawable.Init()
	headerLines := drain.UnitEager(size, header)

	footer := vm.Footer.ToUnit()
	footer.Drawable.Init()
	footerLines := drain.UnitEager(size, footer)

	staticRows := winsize.Rows(
		len(headerLines) + len(footerLines),
	)

	if staticRows > size.Rows {
		return uiState, []text.Line{
			*text.NewLine("Too low resolution"),
		}
	}

	ctx := newRenderContext()

	renderer := pagerRenderer(uiState, *vm.Pager, ctx)

	kernel := vm.Kernel.
		SetRenderer(renderer).
		ToUnit()

	kernel.Drawable.Init()

	dynamicSize := winsize.New(
		size.Rows.Sub(staticRows),
		size.Cols,
	)

	kernelLines, _ := kernel.Drawable.Draw(dynamicSize)
	uiState = syncUIState(uiState, ctx)

	lines := headerLines
	lines = append(lines, kernelLines...)
	lines = append(lines, footerLines...)

	return uiState, lines
}

func syncUIState(uiState *state.UIState, ctx *renderContext) *state.UIState {
	uiState.Pager.ConfirmPage(ctx.MaxPage)
	uiState.Pager.HasMore = ctx.HasMore
	return uiState
}
