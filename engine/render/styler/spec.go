package styler

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/helper"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

type SpecStyler func(style.Spec, winsize.Cols, string, winsize.Cols) (string, bool)

func ps(k style.SpecKind, s SpecStyler) dict.Pair[style.SpecKind, SpecStyler] {
	return dict.NewPair(k, s)
}

var Specs = dict.NewInmutableLinkedMap(
	ps(style.SpcKindFill, func(spec style.Spec, cols winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return fill(spec, cols, text, logicalSize), true
	}),
	ps(style.SpcKindTrimLeft, func(spec style.Spec, _ winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return trimLeft(spec, text, logicalSize), false
	}),
	ps(style.SpcKindTrimRight, func(spec style.Spec, _ winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return trimRight(spec, text, logicalSize), false
	}),
	ps(style.SpcKindPaddingCenter, func(spec style.Spec, cols winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return paddingCenter(spec, cols, text, logicalSize), false
	}),
	ps(style.SpcKindPaddingLeft, func(spec style.Spec, cols winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return paddingLeft(spec, cols, text, logicalSize), false
	}),
	ps(style.SpcKindPaddingRight, func(spec style.Spec, cols winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return paddingRight(spec, cols, text, logicalSize), false
	}),
	ps(style.SpcKindRepeatLeft, func(spec style.Spec, cols winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return repeatLeft(spec, cols, text, logicalSize), false
	}),
	ps(style.SpcKindRepeatRight, func(spec style.Spec, cols winsize.Cols, text string, logicalSize winsize.Cols) (string, bool) {
		return repeatRight(spec, cols, text, logicalSize), false
	}),
)

type Spec struct {
	table *dict.LinkedMap[style.SpecKind, SpecStyler]
}

func NewSpec() *Spec {
	instance := &Spec{}
	return instance.lazyInit()
}

func NewDefaultSpec() *Spec {
	return &Spec{
		table: Specs.Clone(false),
	}
}

func (a *Spec) lazyInit() *Spec {
	if a.table != nil {
		return a
	}

	a.table = dict.NewLinkedMap[style.SpecKind, SpecStyler]()
	return a
}

func (a *Spec) Push(pair ...dict.Pair[style.SpecKind, SpecStyler]) *Spec {
	a.lazyInit()

	a.table.SetPairs(pair...)
	return a
}

func (a *Spec) Apply(spec style.Spec, size winsize.Winsize, text string, logicalSize winsize.Cols) string {
	a.lazyInit()

	exit := false
	kind := spec.Kind()
	for k, p := range a.table.All() {
		if !kind.HasAny(k) {
			continue
		}

		text, exit = p(spec, size.Cols, text, logicalSize)
		if exit {
			return text
		}

		logicalSize = style.SpecMeasureOf(k, spec, style.LayoutContext{
			Cols:     size.Cols,
			TextSize: logicalSize,
		})
	}

	return text
}

func fill(styl style.Spec, cols winsize.Cols, data string, logicalSize winsize.Cols) string {
	opts := helper.LogicalSizeOpts{
		LogicalSize: logicalSize,
	}

	args := styl.Args()

	size := commons.Mapd(args[style.KeyFillSize], cols)
	size = min(cols, size)

	return helper.FillRightWithOpts(data, size, opts)
}

func trimLeft(styl style.Spec, data string, logicalSize winsize.Cols) string {
	if data == "" {
		return data
	}

	args := styl.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyTrimLeftSize], 0)
	size = max(1, size)

	elip := args[style.KeyTrimEllipsisText].Stringf()

	opts := helper.TextTrimOpts{
		LogicalSize:  logicalSize,
		EllipsisText: elip,
	}

	return helper.TrimLeft(data, size, opts)
}

func trimRight(styl style.Spec, data string, logicalSize winsize.Cols) string {
	if data == "" {
		return data
	}

	args := styl.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyTrimRightSize], 0)
	size = max(1, size)

	elip := args[style.KeyTrimEllipsisText].Stringf()

	opts := helper.TextTrimOpts{
		LogicalSize:  logicalSize,
		EllipsisText: elip,
	}

	return helper.TrimRight(data, size, opts)
}

// TODO: Explore the risks of using cols as default
func paddingCenter(styl style.Spec, cols winsize.Cols, data string, logicalSize winsize.Cols) string {
	args := styl.Args()

	size := commons.Mapd(args[style.KeyPaddingCenterSize], cols)
	size = min(cols, size)

	text := args[style.KeyPaddingCenterText].Stringf()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Text:        text,
	}

	return helper.CenterWithOpts(data, size, opts)
}

func paddingLeft(styl style.Spec, cols winsize.Cols, data string, logicalSize winsize.Cols) string {
	args := styl.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyPaddingLeftSize], 0)
	size = min(cols, size)

	text := args[style.KeyPaddingLeftText].Stringf()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Text:        text,
	}

	return helper.LeftWithOpts(data, size, opts)
}

func paddingRight(styl style.Spec, cols winsize.Cols, data string, logicalSize winsize.Cols) string {
	args := styl.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyPaddingRightSize], 0)
	size = min(cols, size)

	text := args[style.KeyPaddingRightText].Stringf()

	opts := helper.TextLayoutOpts{
		LogicalSize: logicalSize,
		Text:        text,
	}

	return helper.RightWithOpts(data, size, opts)
}

func repeatLeft(styl style.Spec, cols winsize.Cols, data string, logicalSize winsize.Cols) string {
	args := styl.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyRepeatLeftSize], 0)
	text := args[style.KeyRepeatLeftText].Stringf()

	if text == "" {
		text = data
		data = ""
	}

	opts := helper.LogicalSizeOpts{
		LogicalSize: logicalSize,
	}

	return helper.RepeatLeftWithOpts(data, text, min(cols, size), opts)
}

func repeatRight(styl style.Spec, cols winsize.Cols, data string, logicalSize winsize.Cols) string {
	args := styl.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyRepeatRightSize], 0)
	text := args[style.KeyRepeatRightText].Stringf()

	if text == "" {
		text = data
		data = ""
	}

	opts := helper.LogicalSizeOpts{
		LogicalSize: logicalSize,
	}

	return helper.RepeatRightWithOpts(data, text, min(cols, size), opts)
}
