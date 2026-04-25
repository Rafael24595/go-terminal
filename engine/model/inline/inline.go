package inline

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

const DefaultInlineSeparator = " | "

type Target uint8

const (
	TargetCode Target = iota
	TargetTags
)

type FilterMeta struct {
	Target Target
	Values set.Set[string]
}

func NewFilterMeta(target Target, values ...string) FilterMeta {
	return FilterMeta{
		Target: target,
		Values: set.SetFrom(values...),
	}
}
