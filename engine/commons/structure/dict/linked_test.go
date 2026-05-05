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

	assert.Panic(t, func() {
		other := NewLinkedMap[string, int]()
		other.SetPairs([]Pair[string, int]{
			{Key: "golang", Value: 1},
			{Key: "rust", Value: 2},
		}...)

		m.Merge(other)
	})

	assert.Panic(t, func() {
		m.SetPairs(Pair[string, int]{"B", 2})
	})

	mc := m.Clone()
	assert.True(t, mc.inmu)
}

func TestLinkedMap_MergeOverride(t *testing.T) {
	m1 := NewLinkedMap[string, int]()
	m1.Set("A", 1)
	m1.Set("B", 2)

	m2 := NewLinkedMap[string, int]()
	m2.Set("B", 20)
	m2.Set("C", 30)

	added, ok := m1.Merge(m2)
	
	assert.True(t, ok)
	assert.Equal(t, 1, added)

	val, _ := m1.Get("B")
	assert.Equal(t, 20, val)
}

func TestLinkedMap_MergeOrder(t *testing.T) {
	m1 := NewLinkedMap[string, int]()
	m1.Set("A", 1)

	m2 := NewLinkedMap[string, int]()
	m2.Set("B", 2)
	m2.Set("C", 3)

	m1.Merge(m2)

	expected := []string{"A", "B", "C"}

	i := 0
	for k := range m1.Keys() {
		assert.Equal(t, expected[i], k)
		i += 1
	}
}

func TestLinkedMap_SetPairsInsertion(t *testing.T) {
	m := NewLinkedMap[string, int]()
	added, ok := m.SetPairs(
		Pair[string, int]{"A", 1},
		Pair[string, int]{"B", 2},
		Pair[string, int]{"C", 3},
	)

	assert.True(t, ok)
	assert.Equal(t, 3, added)
	assert.Equal(t, 3, len(m.data))
}

func TestLinkedMap_SetPairsOverride(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Set("A", 1)

	added, ok := m.SetPairs(Pair[string, int]{"A", 10}, Pair[string, int]{"B", 20})

	assert.True(t, ok)
	assert.Equal(t, 1, added)

	v, _ := m.Get("A")
	assert.Equal(t, 10, v)
}

func TestLinkedMap_CloneIntegrity(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Set("A", 1)
	m.Set("B", 2)

	mc := m.Clone()

	vc, _ := mc.Get("A")
	assert.Equal(t, 1, vc)

	mc.Set("A", 999)
	mc.Set("C", 3)

	v, _ := m.Get("A")

	assert.Equal(t, 1, v)
	assert.False(t, m.Exists("C"))
}

func TestLinkedMap_CloneOrder(t *testing.T) {
	m := NewLinkedMap[int, string]()
	m.Set(1, "A")
	m.Set(2, "B")
	m.Set(3, "C")

	mc := m.Clone()

	keys := []int{}
	for k := range m.Keys() {
		keys = append(keys, k)
	}

	i := 0
	for k := range mc.Keys() {
		assert.Equal(t, keys[i], k)
		i += 1
	}
}
