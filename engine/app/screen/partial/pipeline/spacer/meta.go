package spacer

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"

type Insertion uint8

const (
	Once Insertion = iota
	Between
)

type Meta struct {
	Size      uint8
	Insertion Insertion
	Position  pipeline.Placement
}

func NewMeta(size uint8, insertion Insertion, placement pipeline.Placement) Meta {
	return Meta{
		Size:      size,
		Insertion: insertion,
		Position:  placement,
	}
}
