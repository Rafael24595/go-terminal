package screen

type ScreenArgument[T any] string

func (a ScreenArgument[T]) Type() T {
	var zero T
	return zero
}

func (a ScreenArgument[T]) Code() string {
	return string(a)
}
