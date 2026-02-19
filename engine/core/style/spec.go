package style

import (
	"maps"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/commons"
)

type argMap = map[SpcArgKey]commons.Argument

type SpecsKind uint64

const (
	SpcKindNone SpecsKind = 0

	SpcKindPaddingLeft SpecsKind = 1 << iota
	SpcKindPaddingRight
	SpcKindPaddingCenter

	SpcKindRepeatLeft
	SpcKindRepeatRight

	SpcKindFill
	SpcKindFillUp
	SpcKindFillDown
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

func specDirection(
	kind SpecsKind,
	sizeKey,
	textKey SpcArgKey,
	size uint,
	text ...string,
) Spec {
	args := make(argMap)

	if size > 0 {
		args[sizeKey] = *commons.ArgumentFrom(size)
	}

	if len(text) > 0 {
		args[textKey] = *commons.ArgumentFrom(strings.Join(text, ""))
	}

	return Spec{
		kind: kind,
		args: args,
	}
}
