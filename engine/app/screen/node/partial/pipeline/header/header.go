package header

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/builder"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "header_transformer"

func Transformer(placement pipeline.Placement, lines ...text.Line) pipeline.Transformer {
	drawable := builder.DrainFromLines(lines...)
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

func Node(node screen.Node, placement pipeline.Placement, lines ...text.Line) screen.Node {
	transformer := Transformer(placement, lines...)
	return pipeline.New(node, transformer).ToNode()
}
