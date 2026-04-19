package list

import (
	"iter"
	"sync"
)

type List[T any] struct {
	init sync.Once
	root Item[T]
	size uint
}

func New[T any]() *List[T] {
	return new(List[T]).Init()
}

func (l *List[T]) lazyInit() *List[T] {
	return l.Init()
}

func (l *List[T]) Init() *List[T] {
	l.init.Do(func() {
		l.root.next = &l.root
		l.root.prev = &l.root
		l.size = 0
	})
	return l
}

func (l *List[T]) Size() uint {
	return l.size
}

func (l *List[T]) First() (*Item[T], bool) {
	l.lazyInit()

	if l.size == 0 {
		return nil, false
	}

	return l.root.next, true
}

func (l *List[T]) Last() (*Item[T], bool) {
	l.lazyInit()

	if l.size == 0 {
		return nil, false
	}

	return l.root.prev, true
}

func (l *List[T]) Unshift(data T) *Item[T] {
	l.lazyInit()

	item := newItem(data)
	return l.insert(item, &l.root)
}

func (l *List[T]) Push(data T) *Item[T] {
	l.lazyInit()

	item := newItem(data)
	return l.insert(item, l.root.prev)
}

func (l *List[T]) insert(it, at *Item[T]) *Item[T] {
	it.prev = at
	it.next = at.next

	it.prev.next = it
	it.next.prev = it

	it.list = l
	l.size += 1

	return it
}

func (l *List[T]) Delete(e *Item[T]) (T, bool) {
	if e == nil || e.list != l {
		var zero T
		return zero, false
	}

	e.prev.next = e.next
	e.next.prev = e.prev

	e.next = nil
	e.prev = nil
	e.list = nil

	if l.size > 0 {
		l.size -= 1
	}

	return e.Data, true
}

func (l *List[T]) All() iter.Seq[*Item[T]] {
	return func(yield func(*Item[T]) bool) {
		for i := l.root.next; i != &l.root; i = i.next {
			if !yield(i) {
				return
			}
		}
	}
}
