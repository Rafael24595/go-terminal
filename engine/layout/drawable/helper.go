package drawable

import "github.com/Rafael24595/go-terminal/engine/render/text"

func MaxLineSize(cols int, lines ...text.Line) uint {
	size := uint(0)
	for _, l := range lines {
		measure := text.FragmentMeasure(cols, l.Text...)
		size = max(size, uint(measure))
	}
	return size
}
