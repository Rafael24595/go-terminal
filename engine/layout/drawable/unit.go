package drawable

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

type Unit struct {
	Name     string
	Tags     set.Set[string]
	Drawable Drawable
}

func (c Unit) AddTag(tags ...string) Unit {
	if c.Tags == nil {
		c.Tags = make(set.Set[string])
	}

	c.Tags.Add(tags...)
	return c
}

func IsUnitZero(unit Unit) bool {
	if unit.Name == "" {
		return true
	}

	if unit.Tags == nil {
		return true
	}

	if IsDrawableZero(unit.Drawable) {
		return true
	}

	return false
}
