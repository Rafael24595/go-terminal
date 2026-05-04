package screen

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

const (
	SystemMetaTag = "system_meta"
)

type ScreenMeta struct {
	code set.Set[string]
	tags set.Set[string]
}

func newMeta() ScreenMeta {
	return ScreenMeta{
		code: set.NewSet[string](),
		tags: set.NewSet[string](),
	}
}
