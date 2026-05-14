package drain

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/utils/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameDrain = "drain_pipeline"

func DrawTransformer(lazy bool) pipeline.DrawTransformer {
	return func(size winsize.Winsize, drw drawable.Drawable) ([]text.Line, bool) {
		return drain.Drawable(size, drw, lazy)
	}
}

func Drawable(drawable drawable.Drawable) drawable.Drawable {
	drw := pipeline.New(drawable).
		SetDrawStep(DrawTransformer(true)).
		ToDrawable()

	drw.Name = NameDrain
	return drw
}

func DrawableFromLines(lines ...text.Line) drawable.Drawable {
	return Drawable(
		line.FromLines(lines...).ToDrawable(),
	)
}

func DrawableFromFragments(frags ...text.Fragment) drawable.Drawable {
	return DrawableFromLines(
		*text.LineFromFragments(frags...),
	)
}

func DrawableFromString(txt ...string) drawable.Drawable {
	return DrawableFromLines(
		*text.LineFromFragments(
			text.FragmentsFromString(txt...)...,
		),
	)
}
