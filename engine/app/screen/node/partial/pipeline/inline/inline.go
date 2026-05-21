package inline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"

	drawable_inline "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/inline"
)

const DefaultSeparator = " | "

const Name = "inline_transformer"

type predicate func(pipeline.Filter, drawable.Unit) bool

var predicates = map[pipeline.Criterion]predicate{
	pipeline.Name: func(f pipeline.Filter, d drawable.Unit) bool {
		return f.Values.Has(d.Name)
	},
	pipeline.Tags: func(f pipeline.Filter, d drawable.Unit) bool {
		return d.Tags.Any(f.Values)
	},
}

func Transformer(
	separator string,
	filter pipeline.Filter,
	section pipeline.Section,
	placement pipeline.Placement,
) pipeline.Transformer {
	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		accessor, ok := pipeline.FindViewModelAccessor(section)
		if !ok {
			assert.Unreachable("unsupported target '%d'", section)
			return vm
		}

		vStack := accessor.Get(vm)
		units := vStack.Units()

		unitsLen := len(units)
		matched := make([]drawable.Unit, 0, unitsLen)
		remaining := make([]drawable.Unit, 0, unitsLen)

		for _, d := range units {
			if matchesFilter(filter, d) {
				matched = append(matched, d)
			} else {
				remaining = append(remaining, d)
			}
		}

		if len(matched) == 0 {
			return vm
		}

		inlineUnit := drawable_inline.New(matched...).
			Separator(separator).
			ToUnit()

		inlineUnit.Name = Name

		newVStack := stack.NewVStack(remaining...)

		switch placement {
		case pipeline.After:
			newVStack.Unshift(inlineUnit)
		case pipeline.Before:
			newVStack.Push(inlineUnit)
		default:
			assert.Unreachable("unhandled placement %d", placement)
		}

		return accessor.Set(vm, newVStack)
	}
}

func matchesFilter(filter pipeline.Filter, unit drawable.Unit) bool {
	predicate, ok := predicates[filter.Criterion]
	if !ok {
		return false
	}
	return predicate(filter, unit)
}
