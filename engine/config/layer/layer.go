package layer

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/config/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
)

type config[T math.Number] struct {
	unit   drawable.Unit
	chunk  chunk.Chunk[T]
	static bool
}

type Layer[T math.Number] struct {
	config config[T]
	Value  T
	Status bool
}

func New[T math.Number](unit drawable.Unit, opts ...Option[T]) Layer[T] {
	return FromLayer(
		defaultConfig[T](unit), opts...,
	)
}

func FromLayer[T math.Number](other Layer[T], opts ...Option[T]) Layer[T] {
	cfg := other
	for _, opt := range opts {
		opt(&cfg)
	}

	assert.LazyFalse(func() bool {
		return drawable.IsUnitZero(cfg.Unit())
	}, "unit is not defined")

	cfg.Status = true

	return cfg
}

func (l Layer[T]) Unit() drawable.Unit {
	return l.config.unit
}

func (l Layer[T]) Chunk() chunk.Chunk[T] {
	return l.config.chunk
}
