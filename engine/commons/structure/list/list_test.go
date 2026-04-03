package list

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestList_LazyInit(t *testing.T) {
	var l List[int]

	l.Push(11)

	assert.Equal(t, uint(1), l.Size())
	assert.NotNil(t, l.All())
}

func TestList_BasicOperations(t *testing.T) {
	l := New[string]()

	l.Push("A")
	l.Push("B")
	l.Push("C")

	assert.Equal(t, uint(3), l.Size())

	l.Unshift("0")

	expected := []string{"0", "A", "B", "C"}

	i := 0
	for item := range l.All() {
		assert.Equal(t, expected[i], item.Data)
		i += 1
	}
}

func TestList_FullCycle(t *testing.T) {
	l := New[string]()

	l.Push("B")
	l.Push("C")
	l.Unshift("A")

	assert.Equal(t, uint(3), l.Size())

	itA, ok := l.First()
	assert.True(t, ok)

	itB, ok := itA.Next()
	assert.True(t, ok)
	assert.Equal(t, "B", itB.Data)

	itC, ok := itB.Next()
	assert.True(t, ok)
	assert.Equal(t, "C", itC.Data)

	_, ok = itC.Next()
	assert.False(t, ok)

	_, ok = itA.Prev()
	assert.False(t, ok)
}

func TestList_Delete(t *testing.T) {
	l := New[int]()
	it1 := l.Push(10)
	it2 := l.Push(20)
	it3 := l.Push(30)

	val, ok := l.Delete(it2)

	assert.True(t, ok)
	assert.Equal(t, 20, val)
	assert.Equal(t, uint(2), l.Size())

	next, _ := it1.Next()
	assert.Equal(t, it3, next)

	prev, _ := it3.Prev()
	assert.Equal(t, it1, prev)

	_, ok = l.Delete(it2)
	assert.False(t, ok)
}

func TestList_CrossListDelete(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	item1 := l1.Push(1)

	val, ok := l2.Delete(item1)

	assert.False(t, ok)
	assert.Equal(t, 0, val)
	assert.Equal(t, uint(1), l1.Size())
}
