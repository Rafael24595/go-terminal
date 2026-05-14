package isolated

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/wipe"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameIsolated = "isolated_pipeline"

func Drawable(drawable drawable.Drawable) drawable.Drawable {
	drw := pipeline.New(drawable).
		PushInitSteps(wipe.InitTransformer()).
		SetDrawStep(drain.DrawTransformer(true)).
		ToDrawable()

	drw.Name = NameIsolated
	return drw
}

func DrawableFromLines(lines ...text.Line) drawable.Drawable {
	return Drawable(
		line.FromLines(lines...).ToDrawable(),
	)
}
