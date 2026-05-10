package composer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/page"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func Standard(
	uiState *state.UIState,
	vm viewmodel.ViewModel,
	size winsize.Winsize,
) (*state.UIState, []text.Line) {
	header, footer := vm.InitStaticLayers()

	headerLines := drain.DrawableEager(size, header)
	footerLines := drain.DrawableEager(size, footer)

	staticRows := winsize.Rows(
		len(headerLines) + len(footerLines),
	)

	if staticRows > size.Rows {
		return uiState, []text.Line{
			*text.NewLine("Too low resolution"),
		}
	}

	kernel := vm.InitDynamicLayers()

	dynamicSize := winsize.New(
		size.Rows.Clamp(staticRows),
		size.Cols,
	)

	renderer := page.NewPageRenderer(*vm.Pager)
	status := renderer(uiState, dynamicSize, kernel)

	uiState = syncUIState(uiState, status)

	lines := headerLines
	lines = append(lines, status.Buffer...)
	lines = append(lines, footerLines...)

	return uiState, lines
}

func syncUIState(uiState *state.UIState, status *draw.DrawState) *state.UIState {
	uiState.Pager.ConfirmPage(status.Page)
	uiState.Pager.HasMore = status.ShowPagination()
	return uiState
}
