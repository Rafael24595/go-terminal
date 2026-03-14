package wrapper

import (
	"testing"

	"github.com/Rafael24595/go-terminal/engine/commons/structure/set"
	"github.com/Rafael24595/go-terminal/engine/core/drawable"
	"github.com/Rafael24595/go-terminal/engine/core/inline"
	"github.com/Rafael24595/go-terminal/engine/core/screen"
	"github.com/Rafael24595/go-terminal/test/support/assert"

	drawable_inline "github.com/Rafael24595/go-terminal/engine/core/drawable/inline"
	screen_test "github.com/Rafael24595/go-terminal/test/engine/core/screen"
)

func TestInline_ToScreen(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "base",
	}

	h := NewInline(mock.ToScreen())
	scrn := h.ToScreen()

	screen_test.Helper_ToScreen(t, scrn)

	assert.Equal(t, scrn.Name(), "base")
}

func TestInline_GroupDrawables_NoMatches(t *testing.T) {
	mock := screen_test.MockScreen{}

	w := NewInline(mock.ToScreen())

	filter := inline.NewFilterMeta(
		inline.TargetCode, "x",
	)

	fn := w.groupDrawables(filter)

	d1 := drawable.Drawable{Code: "a"}
	d2 := drawable.Drawable{Code: "b"}

	result := fn(d1, d2)

	assert.Equal(t, 2, len(result))

	assert.Equal(t, "a", result[0].Code)
	assert.Equal(t, "b", result[1].Code)
}

func TestInline_GroupDrawables_ByCode(t *testing.T) {
	mock := screen_test.MockScreen{}

	w := NewInline(mock.ToScreen())

	filter := inline.NewFilterMeta(
		inline.TargetCode, "a",
	)

	fn := w.groupDrawables(filter)

	d1 := drawable.Drawable{Code: "a"}
	d2 := drawable.Drawable{Code: "b"}

	result := fn(d1, d2)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "b", result[0].Code)
	assert.Equal(t, drawable_inline.NameInlineDrawable, result[1].Name)
}

func TestInline_GroupDrawables_ByTags(t *testing.T) {
	w := NewInline(screen.Screen{})

	filter := inline.NewFilterMeta(
		inline.TargetTags, "tag",
	)

	fn := w.groupDrawables(filter)

	d1 := drawable.Drawable{
		Tags: set.SetFrom("tag"),
	}

	d2 := drawable.Drawable{
		Tags: set.SetFrom("other"),
	}

	result := fn(d1, d2)

	assert.Equal(t, 2, len(result))
}

func TestInline_GroupDrawables_MultipleMatches(t *testing.T) {
	w := NewInline(screen.Screen{})

	filter := inline.NewFilterMeta(
		inline.TargetCode, "a",
	)

	fn := w.groupDrawables(filter)

	d1 := drawable.Drawable{Code: "a"}
	d2 := drawable.Drawable{Code: "a"}
	d3 := drawable.Drawable{Code: "b"}

	result := fn(d1, d2, d3)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "b", result[0].Code)
	assert.Equal(t, drawable_inline.NameInlineDrawable, result[1].Name)
}
