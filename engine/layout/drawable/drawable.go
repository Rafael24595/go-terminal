package drawable

import (
	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/render/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type Drawable struct {
	Name string
	Code string
	Tags set.Set[string]
	Init func()
	Wipe func()
	Draw func(size terminal.Winsize) ([]text.Line, bool)
}

func (d Drawable) SetCode(code string) Drawable {
	d.Code = code
	return d
}

func (d Drawable) AddTag(tags ...string) Drawable {
	if d.Tags == nil {
		d.Tags = make(set.Set[string])
	}

	for _, t := range tags {
		d.Tags[t] = struct{}{}
	}
	return d
}
