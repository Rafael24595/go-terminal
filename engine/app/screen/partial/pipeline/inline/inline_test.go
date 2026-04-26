package inline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/partial/pipeline"
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

func TestInline_GroupDrawables_NoMatches(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Code, "x",
	)

	mc1 := drawable_test.MockDrawable{
		Code: "a",
	}
	mc2 := drawable_test.MockDrawable{
		Code: "b",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mc1.ToDrawable(),
			mc2.ToDrawable(),
		)

		transformer := InlineTransformer("|", meta, v)
		vm = transformer(vm)

		items := acc.Get(vm).Items()

		assert.Equal(t, 2, len(items))
		assert.Equal(t, "a", items[0].Code)
		assert.Equal(t, "b", items[1].Code)
	}
}

func TestInline_GroupDrawables_ByCode(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Code, "a",
	)

	mc1 := drawable_test.MockDrawable{
		Code: "a",
	}
	mc2 := drawable_test.MockDrawable{
		Code: "b",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mc1.ToDrawable(),
			mc2.ToDrawable(),
		)

		transformer := InlineTransformer("|", meta, v)
		vm = transformer(vm)

		items := acc.Get(vm).Items()

		assert.Equal(t, 2, len(items))
		assert.Equal(t, "b", items[0].Code)
		assert.Equal(t, name, items[1].Name)
	}
}

func TestInline_GroupDrawables_ByTags(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Tags, "a",
	)

	mc1 := drawable_test.MockDrawable{
		Tags: set.SetFrom("a"),
	}
	mc2 := drawable_test.MockDrawable{
		Tags: set.SetFrom("b"),
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mc1.ToDrawable(),
			mc2.ToDrawable(),
		)

		transformer := InlineTransformer("|", meta, v)
		vm = transformer(vm)

		items := acc.Get(vm).Items()

		assert.Equal(t, 2, len(items))
		assert.Equal(t, name, items[1].Name)
	}
}

func TestInline_GroupDrawables_MultipleMatches(t *testing.T) {
	meta := pipeline.NewFilter(
		pipeline.Code, "a",
	)

	mc1 := drawable_test.MockDrawable{
		Code: "a",
	}
	mc2 := drawable_test.MockDrawable{
		Code: "b",
	}
	mc3 := drawable_test.MockDrawable{
		Code: "a",
	}

	for _, v := range targets {
		vm := *viewmodel.NewViewModel()

		acc := findAccesor(t, v)
		acc.Get(vm).Push(
			mc1.ToDrawable(),
			mc2.ToDrawable(),
			mc3.ToDrawable(),
		)

		transformer := InlineTransformer("|", meta, v)
		vm = transformer(vm)

		items := acc.Get(vm).Items()

		assert.Equal(t, 2, len(items))
		assert.Equal(t, "b", items[0].Code)
		assert.Equal(t, name, items[1].Name)
	}
}
