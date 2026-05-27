package composer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/page"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type renderContext struct {
	MaxPage uint
	HasMore bool
}

func newRenderContext() *renderContext {
	return &renderContext{}
}

func pagerRenderer(uiState *state.UIState, strategy pager.PagerStrategy, ctx *renderContext) stack.LayerRenderer {
	renderer := page.NewPageRenderer(strategy)

	return func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		status := renderer(uiState, size, unit)

		ctx.MaxPage = max(ctx.MaxPage, status.Page)
		if status.ShowPagination() {
			ctx.HasMore = true
		}

		return status.Buffer, !status.Work.Finished()
	}
}
