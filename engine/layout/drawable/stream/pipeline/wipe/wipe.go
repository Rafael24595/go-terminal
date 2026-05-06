package wipe

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func InitTransformer() pipeline.InitTransformer {
	return func(size winsize.Winsize, drw drawable.Drawable) drawable.Drawable {
		drw.Wipe()
		return drw
	}
}
