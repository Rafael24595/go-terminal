package focus

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/page"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

// TODO: Add flag to manage non focus drawable?
func DrawTransformer(engine pager.Engine) pipeline.DrawTransformer {
	predicate := pager.PredicateFocus()
	strategy := *pager.NewStrategy().
		SetEngine(engine).
		SetPredicate(predicate)

	return func(size winsize.Winsize, drawable drawable.Drawable) ([]text.Line, bool) {
		uiState := state.NewUIState()
		renderer := page.NewPageRenderer(strategy)
		status := renderer(uiState, size, drawable)
		return status.Buffer, false
	}
}
