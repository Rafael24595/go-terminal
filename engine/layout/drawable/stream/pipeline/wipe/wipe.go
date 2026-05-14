package wipe

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameWipe = "wipe_pipeline"

func InitTransformer() pipeline.InitTransformer {
	return func(size winsize.Winsize, drw drawable.Drawable) drawable.Drawable {
		drw.Wipe()
		return drw
	}
}

func DrawTransformer() pipeline.DrawTransformer {
	return func(size winsize.Winsize, drw drawable.Drawable) ([]text.Line, bool) {
		lines, status := drw.Draw(size)
		if !status {
			drw.Wipe()
		}
		return lines, true
	}
}

func Drawable(drawable drawable.Drawable) drawable.Drawable {
	drw := pipeline.New(drawable).
		SetDrawStep(DrawTransformer()).
		ToDrawable()

	drw.Name = NameWipe
	return drw
}

