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
	return func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		return drain.Unit(size, unit, lazy)
	}
}

func Unit(unit drawable.Unit) drawable.Unit {
	unt := pipeline.New(unit).
		SetDrawStep(DrawTransformer(true)).
		ToUnit()

	unt.Name = NameDrain
	return unt
}

func UnitFromLines(lines ...text.Line) drawable.Unit {
	return Unit(
		line.FromLines(lines...).ToUnit(),
	)
}

func UnitFromFragments(frags ...text.Fragment) drawable.Unit {
	return UnitFromLines(
		*text.LineFromFragments(frags...),
	)
}

func UnitFromString(txt ...string) drawable.Unit {
	return UnitFromLines(
		*text.LineFromFragments(
			text.FragmentsFromString(txt...)...,
		),
	)
}
