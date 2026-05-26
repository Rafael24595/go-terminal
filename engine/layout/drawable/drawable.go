package drawable

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type InitFunc func()
type WipeFunc func()
type DrawFunc func(size winsize.Winsize) ([]text.Line, bool)

type Drawable struct {
	Init InitFunc
	Wipe WipeFunc
	Draw DrawFunc
}

func IsDrawableZero(drawable Drawable) bool {
	if drawable.Init == nil {
		return true
	}

	if drawable.Wipe == nil {
		return true
	}

	if drawable.Draw == nil {
		return true
	}

	return false
}
