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

func Unit(unit drawable.Unit) drawable.Unit {
	unt := pipeline.New(unit).
		PushInitSteps(wipe.InitTransformer()).
		SetDrawStep(drain.DrawTransformer(true)).
		ToUnit()

	unt.Name = NameIsolated
	return unt
}

func UnitFromLines(lines ...text.Line) drawable.Unit {
	return Unit(
		line.FromLines(lines...).ToUnit(),
	)
}
