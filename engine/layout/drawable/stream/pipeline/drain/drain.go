package drain

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func DrawTransformer(lazy bool) pipeline.DrawTransformer {
	return func(size winsize.Winsize, drw drawable.Drawable) ([]text.Line, bool) {
		return drawable.DrainDrawable(size, drw, lazy)
	}
}
