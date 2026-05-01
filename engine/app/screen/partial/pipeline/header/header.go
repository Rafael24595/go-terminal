package header

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/block"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "header_transformer"

func HeaderTransformer(placement pipeline.Placement, lines ...text.Line) pipeline.Transformer {
	drawable := block.DrawableFromLines(lines...)
	drawable.Name = Name

	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		switch placement {
		case pipeline.Before:
			vm.Header.Unshift(drawable)
		case pipeline.After:
			vm.Header.Push(drawable)
		}
		return vm
	}
}
