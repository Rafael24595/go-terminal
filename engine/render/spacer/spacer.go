package spacer

type Placement uint8

const (
	Before Placement = iota
	After
)

type Insertion uint8

const (
	Once Insertion = iota
	Between
)

type Meta struct {
	Size      uint8
	Insertion Insertion
	Position  Placement
}

func NewMeta(size uint8, insertion Insertion, placement Placement) Meta {
	return Meta{
		Size:      size,
		Insertion: insertion,
		Position:  placement,
	}
}
