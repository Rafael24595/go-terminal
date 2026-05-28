package dummy

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/template"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

const Name = "dummy"

func ToNode() screen.Node {
	return template.New().
		Name(Name).
		ViewModel(*viewmodel.New()).
		ToNode()
}
