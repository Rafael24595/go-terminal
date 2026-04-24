package style

import (
	"maps"
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
)

type LayoutContext struct {
	Text int
	Cols int
}

type argMap = map[SpcArgKey]commons.Argument

var specMeasureTable = dict.NewInmutableLinkedMap(
	dict.P(SpcKindFill, func(spep Spec, ctx LayoutContext) int {
		return spep.args[KeyFillSize].Intd(ctx.Cols)
	}),
	dict.P(SpcKindTrimLeft, func(spec Spec, ctx LayoutContext) int {
		arg := spec.args[KeyTrimLeftSize].Intd(ctx.Text)
		return min(ctx.Text, arg)
	}),
	dict.P(SpcKindTrimRight, func(spec Spec, ctx LayoutContext) int {
		arg := spec.args[KeyTrimRightSize].Intd(ctx.Text)
		return min(ctx.Text, arg)
	}),
	dict.P(SpcKindPaddingCenter, func(spec Spec, ctx LayoutContext) int {
		arg := spec.args[KeyPaddingCenterSize].Intd(ctx.Cols)
		return min(ctx.Cols, arg)
	}),
	dict.P(SpcKindPaddingLeft, func(spec Spec, ctx LayoutContext) int {
		arg := spec.args[KeyPaddingLeftSize].Intd(ctx.Text)
		return max(ctx.Text, arg)
	}),
	dict.P(SpcKindPaddingRight, func(spec Spec, ctx LayoutContext) int {
		arg := spec.args[KeyPaddingRightSize].Intd(ctx.Text)
		return max(ctx.Text, arg)
	}),
	dict.P(SpcKindRepeatLeft, func(spec Spec, ctx LayoutContext) int {
		arg := spec.args[KeyRepeatLeftSize].Intd(ctx.Text)
		return max(ctx.Text, arg)
	}),
	dict.P(SpcKindRepeatRight, func(spec Spec, ctx LayoutContext) int {
		arg := spec.args[KeyRepeatRightSize].Intd(ctx.Text)
		return max(ctx.Text, arg)
	}),
)

type SpecKind uint64

const (
	SpcKindNone SpecKind = 0

	SpcKindPaddingLeft SpecKind = 1 << iota
	SpcKindPaddingRight
	SpcKindPaddingCenter

	SpcKindRepeatLeft
	SpcKindRepeatRight

	SpcKindTrimLeft
	SpcKindTrimRight

	SpcKindFill
)

var specArgsTable = map[SpecKind][]SpcArgKey{
	SpcKindPaddingLeft: {
		KeyPaddingLeftSize, KeyPaddingLeftText,
	},
	SpcKindPaddingRight: {
		KeyPaddingRightSize, KeyPaddingRightText,
	},
	SpcKindPaddingCenter: {
		KeyPaddingCenterSize, KeyPaddingCenterText,
	},
	SpcKindRepeatLeft: {
		KeyRepeatLeftSize, KeyRepeatLeftText,
	},
	SpcKindRepeatRight: {
		KeyRepeatRightSize, KeyRepeatLeftText,
	},
	SpcKindTrimLeft: {
		KeyTrimLeftSize,
	},
	SpcKindTrimRight: {
		KeyTrimRightSize,
	},
	SpcKindFill: {
		KeyFillSize,
	},
}

func (s SpecKind) HasAny(styles ...SpecKind) bool {
	for _, style := range styles {
		if s&style != 0 {
			return true
		}
	}
	return false
}

func (s SpecKind) HasNone(styles ...SpecKind) bool {
	return !s.HasAny(styles...)
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
	kind SpecKind
	args argMap
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

func EraseSpec(target Spec, styles SpecKind) (Spec, Spec) {
	removedKind := target.kind & styles

	removedSpec := SpecFromKind(removedKind)
	if removedKind == SpcKindNone {
		return target, removedSpec
	}

	for kind, keys := range specArgsTable {
		if removedKind&kind == 0 {
			continue
		}

		for _, key := range keys {
			val, ok := target.args[key]
			if !ok {
				continue
			}

			removedSpec.args[key] = val
			delete(target.args, key)
		}
	}

	target.kind &= ^styles

	return target, removedSpec
}

func (s Spec) Kind() SpecKind {
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

func SpecFromKind(kind SpecKind) Spec {
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

func specSize(kind SpecKind, sizeKey SpcArgKey, size uint) Spec {
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
	kind SpecKind,
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

func SpecMeasureOf(kind SpecKind, spec Spec, ctx LayoutContext) int {
	if p, ok := specMeasureTable.Get(kind); ok {
		return p(spec, ctx)
	}
	return ctx.Text
}

func SpecMeasure(spec Spec, ctx LayoutContext) int {
	for k, p := range specMeasureTable.All() {
		if spec.kind.HasAny(k) {
			ctx.Text = p(spec, ctx)
		}
	}
	return ctx.Text
}
