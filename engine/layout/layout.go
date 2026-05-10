package layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Composer func(*state.UIState, viewmodel.ViewModel, winsize.Winsize) (*state.UIState, []text.Line)

type Layout struct {
	Compose Composer
}

type LayoutBuilder struct {
	transformer *winsize.Transformer
	compose       Composer
}

func NewBuilder(apply Composer) *LayoutBuilder {
	return &LayoutBuilder{
		compose: apply,
	}
}

func (b *LayoutBuilder) Transformer(transformer winsize.Transformer) *LayoutBuilder {
	b.transformer = &transformer
	return b
}

func (b *LayoutBuilder) ToLayout() Layout {
	apply := b.compose
	if b.transformer != nil {
		apply = wrapTransformer(apply, *b.transformer)
	}

	return Layout{
		Compose: apply,
	}
}

func wrapTransformer(compose Composer, transformer winsize.Transformer) Composer {
	return func(state *state.UIState, vm viewmodel.ViewModel, size winsize.Winsize) (*state.UIState, []text.Line) {
		newSize := transformer(size)
		return compose(state, vm, newSize)
	}
}
