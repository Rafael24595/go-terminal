package inline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"

	drawable_inline "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/inline"
)

const DefaultInlineSeparator = " | "

const name = "inline_transformer"

type inlinePredicate func(pipeline.Filter, drawable.Drawable) bool

var inlinePredicates = map[pipeline.Criterion]inlinePredicate{
	pipeline.Code: func(f pipeline.Filter, d drawable.Drawable) bool {
		return f.Values.Has(d.Code)
	},
	pipeline.Tags: func(f pipeline.Filter, d drawable.Drawable) bool {
		return d.Tags.Any(f.Values)
	},
}

func InlineTransformer(separator string, filter pipeline.Filter, section pipeline.Section) pipeline.Transformer {
	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		accessor, ok := pipeline.FindViewModelAccessor(section)
		if !ok {
			assert.Unreachable("unsupported target '%d'", section)
			return vm
		}

		vStack := accessor.Get(vm)
		items := vStack.Items()

		itemsLen := len(items)
		matched := make([]drawable.Drawable, 0, itemsLen)
		remaining := make([]drawable.Drawable, 0, itemsLen)

		for _, d := range items {
			if matchesFilter(filter, d) {
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

func matchesFilter(filter pipeline.Filter, drawable drawable.Drawable) bool {
	predicate, ok := inlinePredicates[filter.Criterion]
	if !ok {
		return false
	}
	return predicate(filter, drawable)
}
