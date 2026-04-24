package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type CheckOption struct {
	Id        string
	Status    bool
	Label     text.Fragment
	Timestamp int64
}

func NewCheckOption(id string, option text.Fragment) CheckOption {
	return CheckOption{
		Id:    id,
		Label: option,
	}
}

func FragmentFromCheckOption(options ...CheckOption) []text.Fragment {
	lines := make([]text.Fragment, len(options))
	for i := range options {
		lines[i] = options[i].Label
	}
	return lines
}
