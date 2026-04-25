package spacer

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/render/spacer"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const name = "spacer_transformer"

type spacerPlacement func(spacer.Meta, drawable.Drawable, *stack.VStackDrawable) *stack.VStackDrawable

func SpacerTransformer(meta spacer.Meta, targets ...pipeline.Target) pipeline.Transformer {
	placeSpacer := resolvePlacement(meta)

	drawable := block.BlockDrawableFromLines(
		buildSpacerLines(meta.Size)...,
	)

	drawable.Name = name

	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		for _, target := range targets {
			accessor, ok := pipeline.FindViewModelAccessor(target)
			if !ok {
				assert.Unreachable("unsupported target '%d'", target)
				continue
			}

			stack := accessor.Get(vm)
			if stack.Size() == 0 {
				continue
			}

			stack = placeSpacer(meta, drawable, stack)
			vm = accessor.Set(vm, stack)
		}
		return vm
	}
}

func resolvePlacement(meta spacer.Meta) spacerPlacement {
	if meta.Position == spacer.Before {
		return prependSpacer
	}
	return appendSpacer
}

func prependSpacer(
	meta spacer.Meta,
	drawable drawable.Drawable,
	vStack *stack.VStackDrawable,
) *stack.VStackDrawable {
	if meta.Insertion == spacer.Once {
		vStack.Unshift(drawable)
		return vStack
	}

	newVStack := stack.NewVStackDrawable()
	for _, h := range vStack.Items() {
		newVStack.Push(drawable, h)
	}

	return newVStack
}

func appendSpacer(
	meta spacer.Meta,
	drawable drawable.Drawable,
	vStack *stack.VStackDrawable,
) *stack.VStackDrawable {
	if meta.Insertion == spacer.Once {
		vStack.Push(drawable)
		return vStack
	}

	newVStack := stack.NewVStackDrawable()
	for _, h := range vStack.Items() {
		newVStack.Push(h, drawable)
	}

	return newVStack
}

func buildSpacerLines(size uint8) []text.Line {
	spaces := make([]text.Line, size)
	for i := range spaces {
		spaces[i] = *text.LineJump()
	}
	return spaces
}
