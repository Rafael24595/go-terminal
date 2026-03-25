package delta

import "github.com/Rafael24595/go-terminal/engine/helper/runes"

type Delta struct {
	Start uint
	End   uint
	Text  string
}

func (d Delta) Measure() uint {
	return runes.Measureu(d.Text)
}

func Apply(buffer []rune, d *Delta) []rune {
	if d.Start > d.End {
		return buffer
	}

	size := uint(len(buffer))
	if d.Start > size || d.End > size {
		return buffer
	}

	runesSize := runes.Measureu(d.Text)

	tail := size - d.End
	total := d.Start + runesSize + tail

	newBuffer := make([]rune, total)

	copy(newBuffer[:d.Start], buffer[:d.Start])
	copy(newBuffer[d.Start:], []rune(d.Text))
	copy(newBuffer[d.Start+runesSize:], buffer[d.End:])

	return newBuffer
}
