package buffer

type Clipboard struct {
	buffer []rune
}

func NewClipboard() *Clipboard {
	return &Clipboard{
		buffer: make([]rune, 0),
	}
}

func (b *Clipboard) Size() uint {
	return uint(len(b.buffer))
}

func (b *Clipboard) Put(rns []rune) *Clipboard {
	b.buffer = make([]rune, len(rns))
	copy(b.buffer, rns)
	return b
}

func (b *Clipboard) Buffer() []rune {
	return b.buffer
}
