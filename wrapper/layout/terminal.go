package wrapper_layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/page"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

// TODO: Implement tokenize lines method to prevent line feed injection.
func TerminalApply(state *state.UIState, vm viewmodel.ViewModel, size winsize.Winsize) []text.Line {
	header, footer := vm.InitStaticLayers()

	headerLines := drain.DrawableEager(size, header)
	footerLines := drain.DrawableEager(size, footer)

	static := winsize.Rows(
		len(headerLines) + len(footerLines),
	)

	if static > size.Rows {
		return []text.Line{
			*text.NewLine("Too low resolution"),
		}
	}

	rest := size.Rows.Clamp(static)
	remSize := winsize.New(rest, size.Cols)
	lines := vm.InitDynamicLayers(remSize)

	renderer := page.NewPageRenderer(*vm.Pager)
	dynamicSize := winsize.New(rest, size.Cols)
	drawStt := renderer(state, dynamicSize, lines)

	state.Pager.ConfirmPage(drawStt.Page)
	state.Pager.HasMore = showPagination(drawStt)

	allLines := headerLines
	allLines = append(allLines, drawStt.Buffer...)
	allLines = append(allLines, footerLines...)

	return allLines
}

func showPagination(stt *draw.DrawState) bool {
	return stt.Page != 0 || stt.HasNext || (stt.Work.HasWorks() && !stt.Work.Finished())
}
