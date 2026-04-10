package sink

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure/dict"
	"github.com/Rafael24595/go-terminal/engine/render/style"
	"github.com/Rafael24595/go-terminal/engine/render/text"
)

var specStylesTable = dict.NewInmutableLinkedMap(
	dict.P(style.SpcKindPaddingLeft, sinkLinePaddingLeft),
	dict.P(style.SpcKindPaddingRight, sinkLinePaddingRight),
	dict.P(style.SpcKindPaddingCenter, sinkLinePaddingCenter),
)

func sinkLinePaddingLeft(spec style.SpecKind, line *text.Line, _ int) *text.Line {
	resSpec, delSpec := style.EraseSpec(line.Spec, spec)
	if delSpec.Kind() == style.SpcKindNone {
		return line
	}

	line.Spec = resSpec
	line.UnshiftFragments(
		*text.EmptyFragment().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingRight(spec style.SpecKind, line *text.Line, _ int) *text.Line {
	resSpec, delSpec := style.EraseSpec(line.Spec, spec)
	if delSpec.Kind() == style.SpcKindNone {
		return line
	}

	line.Spec = resSpec
	line.PushFragments(
		*text.EmptyFragment().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingCenter(spec style.SpecKind, line *text.Line, cols int) *text.Line {
	resSpec, delSpec := style.EraseSpec(line.Spec, spec)
	if delSpec.Kind() == style.SpcKindNone {
		return line
	}

	line.Spec = resSpec

	sze := delSpec.Args()[style.KeyPaddingCenterSize].Intd(cols)
	txt := delSpec.Args()[style.KeyPaddingCenterText].Stringf()

	left := sze / 2
	paddLeft := style.SpecPaddingLeft(uint(left), txt)
	line.UnshiftFragments(
		*text.EmptyFragment().AddSpec(paddLeft),
	)

	right := sze - left
	paddRight := style.SpecPaddingRight(uint(right), txt)
	line.PushFragments(
		*text.EmptyFragment().AddSpec(paddRight),
	)

	return line
}

func ApplySinks(line *text.Line, cols int) *text.Line {
	for k, p := range specStylesTable.All() {
		if !line.Spec.Kind().HasAny(k) {
			continue
		}
		line = p(k, line, cols)
	}
	return line
}
