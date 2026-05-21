package spacer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

var targets = []pipeline.Section{
	pipeline.Footer,
	pipeline.Header,
	pipeline.Kernel,
}

func findAccesor(t *testing.T, s pipeline.Section) pipeline.StackAccessor {
	acc, ok := pipeline.FindViewModelAccessor(s)
	if !ok {
		t.Fatalf("unhandled target %d", s)
	}
	return acc
}

func TestSpacer_AddsHeaderLinesWhenEmpty(t *testing.T) {
	for _, v := range targets {
		meta := NewMeta(1, Once, pipeline.After)
		transformer := Transformer(meta, v)

		vm := transformer(
			*viewmodel.NewViewModel(),
		)

		acc := findAccesor(t, v)
		units := acc.Get(vm).Units()

		assert.Len(t, 0, units)
	}
}

func TestSpacer_AddsHeaderLines(t *testing.T) {
	mc := drawable_test.MockUnit{
		Name: "mock_header",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mc.ToUnit(),
		)

		meta := NewMeta(1, Once, pipeline.After)
		transformer := Transformer(meta, v)

		vm = transformer(vm)

		units := acc.Get(vm).Units()

		assert.Len(t, 2, units)

		assert.Equal(t, "mock_header", units[0].Name)
		assert.Equal(t, Name, units[1].Name)
	}

}

func TestSpacer_HeaderBetween(t *testing.T) {
	mock1 := drawable_test.MockUnit{
		Name: "mock_header_1",
	}
	mock2 := drawable_test.MockUnit{
		Name: "mock_header_2",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mock1.ToUnit(),
			mock2.ToUnit(),
		)

		meta := NewMeta(1, Between, pipeline.After)
		transformer := Transformer(meta, v)

		vm = transformer(vm)

		units := acc.Get(vm).Units()

		assert.Len(t, 4, units)

		assert.Equal(t, "mock_header_1", units[0].Name)
		assert.Equal(t, Name, units[1].Name)

		assert.Equal(t, "mock_header_2", units[2].Name)
		assert.Equal(t, Name, units[3].Name)
	}

}
