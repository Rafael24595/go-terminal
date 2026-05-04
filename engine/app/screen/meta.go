package screen

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

const (
	SystemMetaTag = "system_meta"
)

type Meta struct {
	Code set.Set[string]
	Tags set.Set[string]
}

func NewMeta() Meta {
	return Meta{
		Code: set.NewSet[string](),
		Tags: set.NewSet[string](),
	}
}
