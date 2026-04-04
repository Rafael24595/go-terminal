package list

type Item[T any] struct {
	next *Item[T]
	prev *Item[T]
	list *List[T]
	Data T
}

func newItem[T any](data T) *Item[T] {
	return &Item[T]{
		Data: data,
	}
}

func (e *Item[T]) Next() (*Item[T], bool) {
	if e.list == nil {
		return nil, false
	}

	next := e.next
	if next == &e.list.root {
		return nil, false
	}

	return next, true
}

func (e *Item[T]) Prev() (*Item[T], bool) {
	if e.list == nil {
		return nil, false
	}

	prev := e.prev
	if prev == &e.list.root {
		return nil, false
	}

	return prev, true
}
