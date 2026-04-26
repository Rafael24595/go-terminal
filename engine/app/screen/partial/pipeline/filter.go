package pipeline

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

type Criterion uint8

const (
	Code Criterion = iota
	Tags
)

type Filter struct {
	Criterion Criterion
	Values    set.Set[string]
}

func NewFilter(target Criterion, values ...string) Filter {
	return Filter{
		Criterion: target,
		Values:    set.SetFrom(values...),
	}
}
