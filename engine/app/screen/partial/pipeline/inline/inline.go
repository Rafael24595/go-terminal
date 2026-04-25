package partial

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/model/inline"

	drawable_inline "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/inline"
)

const name = "inline_transformer"

type inlinePredicate func(inline.FilterMeta, drawable.Drawable) bool

var inlinePredicates = map[inline.Target]inlinePredicate{
	inline.TargetCode: func(m inline.FilterMeta, d drawable.Drawable) bool {
		return m.Values.Has(d.Code)
	},
	inline.TargetTags: func(m inline.FilterMeta, d drawable.Drawable) bool {
		return d.Tags.Any(m.Values)
	},
}

func InlineTransformer(separator string, meta inline.FilterMeta, target pipeline.Target) pipeline.Transformer {
	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		accessor, ok := pipeline.FindViewModelAccessor(target)
		if !ok {
			assert.Unreachable("unsupported target '%d'", target)
			return vm
		}

		vStack := accessor.Get(vm)
		items := vStack.Items()

		itemsLen := len(items)
		matched := make([]drawable.Drawable, 0, itemsLen)
		remaining := make([]drawable.Drawable, 0, itemsLen)

		for _, d := range items {
			if matchesFilter(meta, d) {
				matched = append(matched, d)
			} else {
				remaining = append(remaining, d)
			}
		}

		if len(matched) == 0 {
			return vm
		}

		inlineDrawable := drawable_inline.NewInlineDrawable(matched...).
			Separator(separator).
			ToDrawable()

		inlineDrawable.Name = name

		newVStack := stack.NewVStackDrawable(
			append(remaining, inlineDrawable)...,
		)

		return accessor.Set(vm, newVStack)
	}
}

func matchesFilter(meta inline.FilterMeta, drawable drawable.Drawable) bool {
	predicate, ok := inlinePredicates[meta.Target]
	if !ok {
		return false
	}
	return predicate(meta, drawable)
}
