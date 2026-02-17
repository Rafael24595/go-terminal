package structure

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/commons/structure"
	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestMinHeap(t *testing.T) {
	h := structure.NewMinHeap[int]()

	h.Push(20)
	h.Push(10)
	h.Push(30)
	h.Push(5)

	val, ok := h.Peek()

	assert.True(t, ok)
	assert.Equal(t, 5, val)
	assert.Equal(t, 4, h.Len())

	expected := []int{5, 10, 20, 30}
	for _, exp := range expected {
		val, ok := h.Pop()
		assert.True(t, ok)
		assert.Equal(t, exp, val)
	}

	_, ok = h.Pop()
	assert.False(t, ok)
}

func TestMaxHeap(t *testing.T) {
	h := structure.NewMaxHeap[int]()

	h.Push(10)
	h.Push(30)
	h.Push(20)

	val, _ := h.Pop()
	assert.Equal(t, 30, val)
	val, _ = h.Pop()
	assert.Equal(t, 20, val)
}

func TestMaxHeapBy_Struct(t *testing.T) {
	type Task struct {
		Name     string
		Priority int
	}

	h := structure.NewMaxHeapBy(func(t Task) int {
		return t.Priority
	})

	h.Push(Task{"Golang", 1})
	h.Push(Task{"Ziglang", 10})
	h.Push(Task{"Rust", 5})

	top, ok := h.Pop()
	assert.True(t, ok)
	assert.Equal(t, "Ziglang", top.Name)

	next, _ := h.Pop()
	assert.Equal(t, "Rust", next.Name)

	last, _ := h.Pop()
	assert.Equal(t, "Golang", last.Name)
}

func TestHeap_EmptyState(t *testing.T) {
    h := structure.NewMinHeap[string]()
    
    val, ok := h.Peek()
    assert.False(t, ok)
    assert.Equal(t, "", val)

    val, ok = h.Pop()
    assert.False(t, ok)
    assert.Equal(t, "", val)
}
