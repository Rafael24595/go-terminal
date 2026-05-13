package builder

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/transformer/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/transformer/wipe"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameIsolated = "isolated_pipeline"

func IsolatedFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	drw := pipeline.New(drawable).
		PushInitSteps(wipe.InitTransformer()).
		SetDrawStep(drain.DrawTransformer(true)).
		ToDrawable()

	drw.Name = NameIsolated
	return drw
}

func IsolatedFromLines(lines ...text.Line) drawable.Drawable {
	return IsolatedFromDrawable(
		line.FromLines(lines...).ToDrawable(),
	)
}
