package layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Build func(*state.UIState, viewmodel.ViewModel, winsize.Winsize) []text.Line

type Layout struct {
	Apply Build
}

type LayoutBuilder struct {
	transformer *winsize.Transformer
	apply       Build
}

func NewBuilder(apply Build) *LayoutBuilder {
	return &LayoutBuilder{
		apply: apply,
	}
}

func (b *LayoutBuilder) Transformer(transformer winsize.Transformer) *LayoutBuilder {
	b.transformer = &transformer
	return b
}

func (b *LayoutBuilder) ToLayout() Layout {
	apply := b.apply
	if b.transformer != nil {
		apply = wrapTransformer(apply, *b.transformer)
	}

	return Layout{
		Apply: apply,
	}
}

func wrapTransformer(apply Build, transformer winsize.Transformer) Build {
	return func(state *state.UIState, vm viewmodel.ViewModel, size winsize.Winsize) []text.Line {
		newSize := transformer(size)
		return apply(state, vm, newSize)
	}
}
