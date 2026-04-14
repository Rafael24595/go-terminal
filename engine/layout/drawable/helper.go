package drawable

import "github.com/Rafael24595/go-terminal/engine/render/text"

func MaxLineSize(lines ...text.Line) uint {
	size := uint(0)
	for _, l := range lines {
		measure := text.FragmentMeasure(l.Text...)
		size = max(size, uint(measure))
	}
	return size
}
