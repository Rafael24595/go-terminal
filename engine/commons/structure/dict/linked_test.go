package dict

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestLinkedMap_ZeroValue(t *testing.T) {
    var m LinkedMap[string, int]
    
    old, exists := m.Set("A", 1)
    assert.False(t, exists)
    assert.Equal(t, 0, old)
    
    val, ok := m.Get("A")
    assert.True(t, ok)
    assert.Equal(t, 1, val)
}

func TestLinkedMap_Order(t *testing.T) {
	m := NewLinkedMap[string, int]()

	m.Set("A", 1)
	m.Set("B", 2)
	m.Set("C", 3)

	keys := []string{"A", "B", "C"}

	i := 0
	for k := range m.Keys() {
		assert.Equal(t, keys[i], k)
		i += 1
	}
}

func TestLinkedMap_UpdateValue(t *testing.T) {
	m := NewLinkedMap[string, int]()

	m.Set("A", 1)
	m.Set("B", 2)

	old, exists := m.Set("A", 10)

	assert.True(t, exists)
	assert.Equal(t, 1, old)

	for k, v := range m.All() {
		assert.Equal(t, "A", k)
		assert.Equal(t, 10, v)
		break
	}
}

func TestLinkedMap_OrderAndUpdates(t *testing.T) {
    m := NewLinkedMap[string, string]()
    
    m.Set("1", "a.1")
    m.Set("2", "b.1")
    m.Set("3", "c.1")
    
    m.Set("2", "b.2")
    
    expectedKeys := []string{"1", "2", "3"}
    expectedValues := []string{"a.1", "b.2", "c.1"}
    
    i := 0
    for k, v := range m.All() {
        assert.Equal(t, expectedKeys[i], k)
        assert.Equal(t, expectedValues[i], v)
       
		i += 1
    }
}

func TestLinkedMap_Delete(t *testing.T) {
	m := NewLinkedMap[string, int]()

	m.Set("X", 100)
	m.Set("Y", 200)
	m.Set("Z", 300)

	old, ok := m.Delete("Y")
	assert.True(t, ok)
	assert.Equal(t, 200, old)

	expected := []string{"X", "Z"}
	i := 0
	for k := range m.Keys() {
		assert.Equal(t, expected[i], k)
		i += 1
	}
}

func TestLinkedMap_Inmutable(t *testing.T) {
	m := NewInmutableLinkedMap(
		[]Pair[string, int]{
			{Key: "golang", Value: 1},
			{Key: "rust", Value: 2},
		}...,
	)

	val, ok := m.Get("golang")

	assert.True(t, ok)
	assert.Equal(t, 1, val)

	assert.Panic(t, func() {
		m.Set("ziglang", 3)
	})

	assert.Panic(t, func() {
		m.Delete("rust")
	})
}
