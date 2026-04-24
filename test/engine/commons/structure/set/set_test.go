package set_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

func TestSet_Has(t *testing.T) {
	s := set.SetFrom("apple", "banana")

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Element exists",
			input:    "apple",
			expected: true,
		},
		{
			input:    "Element does not exist",
			name:     "orange",
			expected: false,
		},
		{
			input:    "Empty string",
			name:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, s.Has(tt.input))
		})
	}
}

func TestSet_Any(t *testing.T) {
	tests := []struct {
		name     string
		setA     []int
		setB     []int
		expected bool
	}{
		{
			name:     "Direct intersection",
			setA:     []int{1, 2, 3},
			setB:     []int{3, 4, 5},
			expected: true,
		},
		{
			name:     "No intersection",
			setA:     []int{1, 2},
			setB:     []int{3, 4},
			expected: false,
		},
		{
			name:     "One set is empty",
			setA:     []int{1, 2},
			setB:     []int{},
			expected: false,
		},
		{
			name:     "Both sets are empty",
			setA:     []int{},
			setB:     []int{},
			expected: false,
		},
		{
			name:     "Identical sets",
			setA:     []int{10, 20},
			setB:     []int{10, 20},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := set.SetFrom(tt.setA...)
			b := set.SetFrom(tt.setB...)

			assert.Equal(t, tt.expected, a.Any(b))
			assert.Equal(t, tt.expected, b.Any(a))
		})
	}
}

func TestSet_Add(t *testing.T) {
	s := set.NewSet[int](1)
	s.Add(42)

	assert.Len(t, 1, s)
	assert.True(t, s.Has(42))
}

func TestSet_Merge(t *testing.T) {
	s1 := set.SetFrom(1, 2, 3)
	s2 := set.SetFrom(3, 4, 5)

	s1.Merge(s2)

	assert.Len(t, 5, s1)

	assert.True(t, s1.Has(4))
	assert.True(t, s1.Has(5))
}

func BenchmarkSet_Any(b *testing.B) {
	large := set.NewSet[int](1000)
	for i := range 1000 {
		large.Add(i)
	}

	small := set.NewSet[int](2)
	small.Add(999)
	small.Add(2000)

	for b.Loop() {
		large.Any(small)
	}
}
