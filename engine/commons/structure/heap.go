package structure

import "cmp"

type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less:  less,
	}
}

func NewMinHeap[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less:  func(a, b T) bool {
			return a < b
		},
	}
}

func NewMaxHeap[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less:  func(a, b T) bool {
			return a > b
		},
	}
}

func NewMinHeapBy[T any, K cmp.Ordered](get func(i T) K) *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less:  func(a, b T) bool {
			return get(a) < get(b)
		},
	}
}

func NewMaxHeapBy[T any, K cmp.Ordered](get func(i T) K) *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less:  func(a, b T) bool {
			return get(a) > get(b)
		},
	}
}

func (h *Heap[T]) Len() int {
	return len(h.items)
}

func (h *Heap[T]) Push(x T) {
	h.items = append(h.items, x)
	h.up(h.Len() - 1)
}

func (h *Heap[T]) Pop() (T, bool) {
	if h.Len() == 0 {
		var zero T
		return zero, false
	}

	top := h.items[0]
	last := h.pop()

	if h.Len() > 0 {
		h.items[0] = last
		h.down(0)
	}

	return top, true
}

func (h *Heap[T]) Peek() (T, bool) {
	if h.Len() == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

func (h *Heap[T]) pop() T {
	last := len(h.items) - 1

	item := h.items[last]
	h.items = h.items[:last]

	return item
}

func (h *Heap[T]) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *Heap[T]) up(index int) {
	for {
		parent := (index - 1) / 2
		if parent == index || !h.less(h.items[index], h.items[parent]) {
			break
		}

		h.swap(parent, index)
		index = parent
	}
}

func (h *Heap[T]) down(index int) bool {
	size := h.Len()

	cursor := index
	for {
		childLeft := 2*cursor + 1
		if childLeft >= size || childLeft < 0 {
			break
		}

		priorityChild := childLeft

		candidateFix := childLeft + 1
		if candidateFix < size && h.less(h.items[candidateFix], h.items[childLeft]) {
			priorityChild = candidateFix
		}

		if !h.less(h.items[priorityChild], h.items[cursor]) {
			break
		}

		h.swap(cursor, priorityChild)
		cursor = priorityChild
	}

	return cursor > index
}
