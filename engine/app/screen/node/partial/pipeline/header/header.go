package header

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "header_transformer"

func Transformer(placement pipeline.Placement, lines ...text.Line) pipeline.Transformer {
	unit := drain.UnitFromLines(lines...)
	unit.Name = Name

	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		switch placement {
		case pipeline.Before:
			vm.Header.Unshift(unit)
		case pipeline.After:
			vm.Header.Push(unit)
		}
		return vm
	}
}

func Node(node screen.Node, lines ...text.Line) screen.Node {
	transformer := Transformer(pipeline.After, lines...)
	return pipeline.New(node, transformer).ExpireOnNode().ToNode()
}
