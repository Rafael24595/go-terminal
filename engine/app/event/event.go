package event

type Event[T any] struct {
	Value T
}

func NewEvent[T any](value T) Event[T] {
	return Event[T]{
		Value: value,
	}
}
