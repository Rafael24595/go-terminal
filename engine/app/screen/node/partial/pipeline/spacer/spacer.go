package spacer

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/wipe"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_pipeline "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
)

const Name = "spacer_transformer"

type spacerPlacement func(Meta, drawable.Drawable, *stack.VStackDrawable) *stack.VStackDrawable

func SpacerTransformer(meta Meta, sections ...pipeline.Section) pipeline.Transformer {
	placeSpacer := resolvePlacement(meta)

	spacer := makeSpacerDrawable(meta.Size)
	
	drawable := drawable_pipeline.New(spacer).
		PushInitSteps(wipe.InitTransformer()).
		SetDrawStep(drain.DrawTransformer(true)).
		ToDrawable()

	drawable.Name = Name

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

			stack = placeSpacer(meta, drawable, stack)
			vm = accessor.Set(vm, stack)
		}
		return vm
	}
}

func resolvePlacement(meta Meta) spacerPlacement {
	if meta.Position == pipeline.Before {
		return prependSpacer
	}
	return appendSpacer
}

func prependSpacer(
	meta Meta,
	drawable drawable.Drawable,
	vStack *stack.VStackDrawable,
) *stack.VStackDrawable {
	if meta.Insertion == Once {
		vStack.Unshift(drawable)
		return vStack
	}

	newVStack := stack.NewVStack()
	for _, h := range vStack.Items() {
		newVStack.Push(drawable, h)
	}

	return newVStack
}

func appendSpacer(
	meta Meta,
	drawable drawable.Drawable,
	vStack *stack.VStackDrawable,
) *stack.VStackDrawable {
	if meta.Insertion == Once {
		vStack.Push(drawable)
		return vStack
	}

	newVStack := stack.NewVStack()
	for _, h := range vStack.Items() {
		newVStack.Push(h, drawable)
	}

	return newVStack
}

func makeSpacerDrawable(size uint8) drawable.Drawable {
	spaces := make([]text.Line, size)
	for i := range spaces {
		spaces[i] = *text.LineJump()
	}
	return line.New(spaces...).ToDrawable()
}
