package form

type pointer uint8

const (
	pointerPrompt pointer = 1 << iota
	pointerGutter
)

func (m pointer) hasAny(pointers ...pointer) bool {
	for _, mod := range pointers {
		if m&mod != 0 {
			return true
		}
	}
	return false
}

func (m pointer) hasNone(pointers ...pointer) bool {
	return !m.hasAny(pointers...)
}

var pointers = []pointer{
	pointerPrompt,
	pointerGutter,
	pointerPrompt | pointerGutter,
}

func findPointer(cursor uint8) pointer {
	if cursor >= uint8(len(pointers)) {
		return pointers[0]
	}
	return pointers[cursor]
}

func nextPointer(cursor uint8) uint8 {
	return (cursor + 1) % uint8(len(pointers))
}
