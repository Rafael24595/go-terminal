package builder

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/transformer/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameDrain = "drain_pipeline"

func DrainFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	drw := pipeline.New(drawable).
		SetDrawStep(drain.DrawTransformer(true)).
		ToDrawable()

	drw.Name = NameDrain
	return drw
}

func DrainFromLines(lines ...text.Line) drawable.Drawable {
	return DrainFromDrawable(
		line.FromLines(lines...).ToDrawable(),
	)
}

func DrainFromFragments(frags ...text.Fragment) drawable.Drawable {
	return DrainFromLines(
		*text.LineFromFragments(frags...),
	)
}

func DrainFromString(txt ...string) drawable.Drawable {
	return DrainFromLines(
		*text.LineFromFragments(
			text.FragmentsFromString(txt...)...,
		),
	)
}
