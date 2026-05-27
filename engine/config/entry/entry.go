package entry

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Entry struct {
	Node       screen.Node
	Selectable bool
	Opts       []layer.Option[winsize.Rows]
}

func New(node screen.Node, opts ...Option) Entry {
	cfg := defaultEntry(node)
	for _, opt := range opts {
		opt(&cfg)
	}

	assert.LazyFalse(func() bool {
		return screen.IsZeroNode(cfg.Node)
	}, "unit is not defined")

	return cfg
}
