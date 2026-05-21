package spacer

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/isolated"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "spacer_transformer"

type placement func(Meta, drawable.Unit, *stack.VStackUnit) *stack.VStackUnit

func Transformer(meta Meta, sections ...pipeline.Section) pipeline.Transformer {
	spacer := resolvePlacement(meta)

	unit := isolated.UnitFromLines(
		makeLines(meta.Size)...,
	)

	unit.Name = Name

	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		for _, section := range sections {
			accessor, ok := pipeline.FindViewModelAccessor(section)
			if !ok {
				assert.Unreachable("unsupported target '%d'", section)
				continue
			}

			stack := accessor.Get(vm)
			if stack.Size() == 0 {
				continue
			}

			stack = spacer(meta, unit, stack)
			vm = accessor.Set(vm, stack)
		}
		return vm
	}
}

func resolvePlacement(meta Meta) placement {
	if meta.Position == pipeline.Before {
		return prependSpacer
	}
	return appendSpacer
}

func prependSpacer(
	meta Meta,
	unit drawable.Unit,
	vStack *stack.VStackUnit,
) *stack.VStackUnit {
	if meta.Insertion == Once {
		vStack.Unshift(unit)
		return vStack
	}

	newVStack := stack.NewVStack()
	for _, h := range vStack.Units() {
		newVStack.Push(unit, h)
	}

	return newVStack
}

func appendSpacer(
	meta Meta,
	unit drawable.Unit,
	vStack *stack.VStackUnit,
) *stack.VStackUnit {
	if meta.Insertion == Once {
		vStack.Push(unit)
		return vStack
	}

	newVStack := stack.NewVStack()
	for _, h := range vStack.Units() {
		newVStack.Push(h, unit)
	}

	return newVStack
}

func makeLines(size uint8) []text.Line {
	spaces := make([]text.Line, size)
	for i := range spaces {
		spaces[i] = *text.LineJump()
	}
	return spaces
}
