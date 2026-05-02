package drawable

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func MaxLineSize(cols winsize.Cols, lines ...text.Line) winsize.Cols {
	size := winsize.Cols(0)
	for _, l := range lines {
		measure := text.FragmentMeasure(cols, l.Text...)
		size = max(size, measure)
	}
	return size
}
