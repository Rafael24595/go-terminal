package style

import (
	"maps"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/commons"
	"github.com/Rafael24595/go-terminal/engine/commons/structure/dict"
)

type LayoutContext struct {
	Cols int
}

type argMap = map[SpcArgKey]commons.Argument

var specMeasureTableWithContext = dict.NewInmutableLinkedMap(
	dict.P(SpcKindFill, func(spep Spec, size int, ctx LayoutContext) int {
		return spep.args[KeyFillSize].Intd(ctx.Cols)
	}),
)

var specMeasureTable = dict.NewInmutableLinkedMap(
	dict.P(SpcKindTrimLeft, func(spec Spec, size int) int {
		arg := spec.args[KeyTrimLeftSize].Intd(size)
		return min(size, arg)
	}),
	dict.P(SpcKindTrimRight, func(spec Spec, size int) int {
		arg := spec.args[KeyTrimRightSize].Intd(size)
		return min(size, arg)
	}),
	dict.P(SpcKindPaddingCenter, func(spec Spec, size int) int {
		arg := spec.args[KeyPaddingCenterSize].Intd(size)
		return max(size, arg)
	}),
	dict.P(SpcKindPaddingLeft, func(spec Spec, size int) int {
		arg := spec.args[KeyPaddingLeftSize].Intd(size)
		return max(size, arg)
	}),
	dict.P(SpcKindPaddingRight, func(spec Spec, size int) int {
		arg := spec.args[KeyPaddingRightSize].Intd(size)
		return max(size, arg)
	}),
	dict.P(SpcKindRepeatLeft, func(spec Spec, size int) int {
		arg := spec.args[KeyRepeatLeftSize].Intd(size)
		return max(size, arg)
	}),
	dict.P(SpcKindRepeatRight, func(spec Spec, size int) int {
		arg := spec.args[KeyRepeatRightSize].Intd(size)
		return max(size, arg)
	}),
)

type SpecsKind uint64

const (
	SpcKindNone SpecsKind = 0

	SpcKindPaddingLeft SpecsKind = 1 << iota
	SpcKindPaddingRight
	SpcKindPaddingCenter

	SpcKindRepeatLeft
	SpcKindRepeatRight

	SpcKindTrimLeft
	SpcKindTrimRight

	SpcKindFill
)

func (s SpecsKind) HasAny(styles ...SpecsKind) bool {
	for _, style := range styles {
		if s&style != 0 {
			return true
		}
	}
	return false
}

func (s SpecsKind) HasNone(styles ...SpecsKind) bool {
	return !s.HasAny(styles...)
}

func MergeSpec(styles ...Spec) Spec {
	kind := SpcKindNone
	args := make(argMap)

	for _, style := range styles {
		kind |= style.kind
		maps.Copy(args, style.args)
	}

	return Spec{
		kind: kind,
		args: args,
	}
}

func EraseSpec(target Spec, styles SpecsKind) Spec {
	target.kind &= ^styles
	//TODO: Delete args.
	return target
}

type SpcArgKey uint8

const (
	KeyPaddingLeftSize SpcArgKey = iota
	KeyPaddingLeftText

	KeyPaddingRightSize
	KeyPaddingRightText

	KeyPaddingCenterSize
	KeyPaddingCenterText

	KeyRepeatLeftSize
	KeyRepeatLeftText

	KeyRepeatRightSize
	KeyRepeatRightText

	KeyTrimLeftSize
	KeyTrimRightSize
	KeyTrimEllipsisText

	KeyFillSize
)

type Spec struct {
	kind SpecsKind
	args argMap
}

func (s Spec) Kind() SpecsKind {
	return s.kind
}

func (s Spec) Args() argMap {
	return s.args
}

func SpecEmpty() Spec {
	return Spec{
		kind: SpcKindNone,
		args: make(argMap),
	}
}

func SpecFromKind(kind SpecsKind) Spec {
	return Spec{
		kind: kind,
		args: make(argMap),
	}
}

func SpecPaddingLeft(size uint, text ...string) Spec {
	return specDirection(
		SpcKindPaddingLeft,
		KeyPaddingLeftSize,
		KeyPaddingLeftText,
		size,
		text...,
	)
}

func SpecPaddingRight(size uint, text ...string) Spec {
	return specDirection(
		SpcKindPaddingRight,
		KeyPaddingRightSize,
		KeyPaddingRightText,
		size,
		text...,
	)
}

func SpecPaddingCenter(size uint, text ...string) Spec {
	return specDirection(
		SpcKindPaddingCenter,
		KeyPaddingCenterSize,
		KeyPaddingCenterText,
		size,
		text...,
	)
}

func SpecRepeatLeft(size uint, text ...string) Spec {
	return specDirection(
		SpcKindRepeatLeft,
		KeyRepeatLeftSize,
		KeyRepeatLeftText,
		size,
		text...,
	)
}

func SpecRepeatRight(size uint, text ...string) Spec {
	return specDirection(
		SpcKindRepeatRight,
		KeyRepeatRightSize,
		KeyRepeatRightText,
		size,
		text...,
	)
}

func SpecTrimLeft(size uint) Spec {
	return specSize(
		SpcKindTrimLeft,
		KeyTrimLeftSize,
		size,
	)
}

func SpecTrimRight(size uint) Spec {
	return specSize(
		SpcKindTrimRight,
		KeyTrimRightSize,
		size,
	)
}

func SpecTrimTextLeft(size uint, ellipsis string) Spec {
	spec := specSize(
		SpcKindTrimLeft,
		KeyTrimLeftSize,
		size,
	)

	spec.args[KeyTrimEllipsisText] = commons.ArgumentFrom(ellipsis)

	return spec
}

func SpecTrimTextRight(size uint, ellipsis string) Spec {
	spec := specSize(
		SpcKindTrimRight,
		KeyTrimRightSize,
		size,
	)

	spec.args[KeyTrimEllipsisText] = commons.ArgumentFrom(ellipsis)

	return spec
}

func SpecFill(size uint) Spec {
	return specSize(
		SpcKindFill,
		KeyFillSize,
		size,
	)
}

func specSize(kind SpecsKind, sizeKey SpcArgKey, size uint) Spec {
	args := make(argMap)

	if size > 0 {
		args[sizeKey] = commons.ArgumentFrom(size)
	}

	return Spec{
		kind: kind,
		args: args,
	}
}

func specDirection(
	kind SpecsKind,
	sizeKey,
	textKey SpcArgKey,
	size uint,
	text ...string,
) Spec {
	args := make(argMap)

	if size > 0 {
		args[sizeKey] = commons.ArgumentFrom(size)
	}

	if len(text) > 0 {
		args[textKey] = commons.ArgumentFrom(strings.Join(text, ""))
	}

	return Spec{
		kind: kind,
		args: args,
	}
}

func SpecMeasureOf(kind SpecsKind, spec Spec, size int) int {
	if predicate, ok := specMeasureTable.Get(kind); ok {
		return predicate(spec, size)
	}
	return size
}

func SpecMeasureWithContext(spec Spec, size int, ctx LayoutContext) int {
	for k, f := range specMeasureTableWithContext.All() {
		if spec.kind.HasAny(k) {
			size = f(spec, size, ctx)
		}
	}
	return SpecMeasure(spec, size)
}

func SpecMeasure(spec Spec, size int) int {
	for k, p := range specMeasureTable.All() {
		if spec.kind.HasAny(k) {
			size = p(spec, size)
		}
	}
	return size
}
