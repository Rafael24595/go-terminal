package dict

import (
	"iter"
	"sync"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/list"
)

type LinkedMap[K comparable, V any] struct {
	init sync.Once
	inmu bool
	list *list.List[Pair[K, V]]
	data map[K]*list.Item[Pair[K, V]]
}

func NewLinkedMap[K comparable, V any]() *LinkedMap[K, V] {
	return new(LinkedMap[K, V]).Init()
}

func NewInmutableLinkedMap[K comparable, V any](pairs ...Pair[K, V]) *LinkedMap[K, V] {
	linked := NewLinkedMap[K, V]().Init()
	linked.inmu = true

	for _, p := range pairs {
		linked.set(p)
	}

	return linked
}

func (m *LinkedMap[K, V]) lazyInit() *LinkedMap[K, V] {
	return m.Init()
}

func (m *LinkedMap[K, V]) Init() *LinkedMap[K, V] {
	m.init.Do(func() {
		m.list = list.New[Pair[K, V]]()
		m.data = make(map[K]*list.Item[Pair[K, V]])
	})
	return m
}

func (m *LinkedMap[K, V]) Exists(k K) bool {
	_, exists := m.Get(k)
	return exists
}

func (m *LinkedMap[K, V]) Get(k K) (V, bool) {
	m.lazyInit()

	if item, exists := m.data[k]; exists {
		return item.Data.Value, true
	}

	var zero V
	return zero, false
}

func (m *LinkedMap[K, V]) Set(k K, v V) (V, bool) {
	var old V

	if m.inmu {
		assert.Unreachable("cannot modify an inmutable souce")
		return old, false
	}

	m.lazyInit()

	pair := NewPair(k, v)

	return m.set(pair)
}

func (m *LinkedMap[K, V]) set(pair Pair[K, V]) (V, bool) {
	var old V

	item, exists := m.data[pair.Key]
	if !exists {
		m.data[pair.Key] = m.list.Push(pair)
		return old, false
	}

	old = item.Data.Value
	item.Data = pair

	return old, true
}

func (m *LinkedMap[K, V]) Delete(k K) (V, bool) {
	var old V

	if m.inmu {
		assert.Unreachable("cannot modify an inmutable souce")
		return old, false
	}

	m.lazyInit()

	item, exists := m.data[k]
	if !exists {
		return old, false
	}

	old = item.Data.Value

	m.list.Delete(item)
	delete(m.data, k)

	return old, true
}

func (m *LinkedMap[K, V]) All() iter.Seq2[K, V] {
	m.lazyInit()

	return func(yield func(K, V) bool) {
		for item := range m.list.All() {
			if !yield(item.Data.Key, item.Data.Value) {
				return
			}
		}
	}
}

func (m *LinkedMap[K, V]) Keys() iter.Seq[K] {
	m.lazyInit()

	return func(yield func(K) bool) {
		for item := range m.list.All() {
			if !yield(item.Data.Key) {
				return
			}
		}
	}
}

func (m *LinkedMap[K, V]) Values() iter.Seq[V] {
	m.lazyInit()

	return func(yield func(V) bool) {
		for item := range m.list.All() {
			if !yield(item.Data.Value) {
				return
			}
		}
	}
}
