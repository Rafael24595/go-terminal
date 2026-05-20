package wipe

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameWipe = "wipe_pipeline"

func InitTransformer() pipeline.InitTransformer {
	return func(size winsize.Winsize, drawable drawable.Drawable) drawable.Drawable {
		drawable.Wipe()
		return drawable
	}
}

func DrawTransformer() pipeline.DrawTransformer {
	transformer := DataTransformer()
	return func(size winsize.Winsize, drawable drawable.Drawable) ([]text.Line, bool) {
		lines, hasNext := drawable.Draw(size)
		return transformer(size, drawable, lines, hasNext)
	}
}

func DataTransformer() pipeline.DataTransformer {
	return func(_ winsize.Winsize, drawable drawable.Drawable, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		if len(lines) == 0 {
			return lines, hasNext
		}

		if !hasNext {
			drawable.Wipe()
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

