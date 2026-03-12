package drawable

import (
	"github.com/Rafael24595/go-terminal/engine/core/text"
	"github.com/Rafael24595/go-terminal/engine/terminal"
)

type TagSet map[string]struct{}

type Drawable struct {
	Code string
	Tags TagSet
	Init func(size terminal.Winsize)
	Draw func() ([]text.Line, bool)
}

func (d Drawable) SetCode(code string) Drawable {
	d.Code = code
	return d
}

func (d Drawable) AddTag(tags ...string) Drawable {
	if d.Tags == nil {
		d.Tags = make(TagSet)	
	}

	for _, t := range tags {
		d.Tags[t] = struct{}{}
	}
	return d
}
