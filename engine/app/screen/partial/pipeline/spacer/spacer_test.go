package spacer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
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
		meta := NewMeta(1, Once, After)
		transformer := SpacerTransformer(meta, v)

		vm := transformer(
			*viewmodel.NewViewModel(),
		)

		acc := findAccesor(t, v)
		items := acc.Get(vm).Items()

		assert.Len(t, 0, items)
	}
}

func TestSpacer_AddsHeaderLines(t *testing.T) {
	mc := drawable_test.MockDrawable{
		Name: "mock_header",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mc.ToDrawable(),
		)

		meta := NewMeta(1, Once, After)
		transformer := SpacerTransformer(meta, v)

		vm = transformer(vm)

		items := acc.Get(vm).Items()

		assert.Len(t, 2, items)

		assert.Equal(t, "mock_header", items[0].Name)
		assert.Equal(t, name, items[1].Name)
	}

}

func TestSpacer_HeaderBetween(t *testing.T) {
	mc1 := drawable_test.MockDrawable{
		Name: "mock_header_1",
	}
	mc2 := drawable_test.MockDrawable{
		Name: "mock_header_2",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mc1.ToDrawable(),
			mc2.ToDrawable(),
		)

		meta := NewMeta(1, Between, After)
		transformer := SpacerTransformer(meta, v)

		vm = transformer(vm)

		items := acc.Get(vm).Items()

		assert.Len(t, 4, items)

		assert.Equal(t, "mock_header_1", items[0].Name)
		assert.Equal(t, name, items[1].Name)

		assert.Equal(t, "mock_header_2", items[2].Name)
		assert.Equal(t, name, items[3].Name)
	}

}
