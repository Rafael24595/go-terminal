package inline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

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

func TestInline_GroupUnits_NoMatches(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Name, "x",
	)

	mock1 := drawable_test.MockUnit{
		Name: "a",
	}
	mock2 := drawable_test.MockUnit{
		Name: "b",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mock1.ToUnit(),
			mock2.ToUnit(),
		)

		transformer := Transformer("|", meta, v, pipeline.After)
		vm = transformer(vm)

		units := acc.Get(vm).Units()

		assert.Equal(t, 2, len(units))
		assert.Equal(t, "a", units[0].Name)
		assert.Equal(t, "b", units[1].Name)
	}
}

func TestInline_GroupUnits_ByCode(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Name, "a",
	)

	mock1 := drawable_test.MockUnit{
		Name: "a",
	}
	mock2 := drawable_test.MockUnit{
		Name: "b",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mock1.ToUnit(),
			mock2.ToUnit(),
		)

		transformer := Transformer("|", meta, v, pipeline.Before)
		vm = transformer(vm)

		units := acc.Get(vm).Units()

		assert.Equal(t, 2, len(units))
		assert.Equal(t, "b", units[0].Name)
		assert.Equal(t, Name, units[1].Name)
	}
}

func TestInline_GroupUnits_ByTags(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Tags, "a",
	)

	mock1 := drawable_test.MockUnit{
		Tags: set.SetFrom("a"),
	}
	mock2 := drawable_test.MockUnit{
		Tags: set.SetFrom("b"),
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mock1.ToUnit(),
			mock2.ToUnit(),
		)

		transformer := Transformer("|", meta, v, pipeline.Before)
		vm = transformer(vm)

		units := acc.Get(vm).Units()

		assert.Equal(t, 2, len(units))
		assert.Equal(t, Name, units[1].Name)
	}
}

func TestInline_GroupUnits_MultipleMatches(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Name, "a",
	)

	mock1 := drawable_test.MockUnit{
		Name: "a",
	}
	mock2 := drawable_test.MockUnit{
		Name: "b",
	}
	mock3 := drawable_test.MockUnit{
		Name: "a",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mock1.ToUnit(),
			mock2.ToUnit(),
			mock3.ToUnit(),
		)

		transformer := Transformer("|", meta, v, pipeline.After)
		vm = transformer(vm)

		units := acc.Get(vm).Units()

		assert.Equal(t, 2, len(units))
		assert.Equal(t, Name, units[0].Name)
		assert.Equal(t, "b", units[1].Name)
	}
}
