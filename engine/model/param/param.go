package param

type Typed[T any] string

func (a Typed[T]) Type() T {
	var zero T
	return zero
}

func (a Typed[T]) Code() string {
	return string(a)
}
