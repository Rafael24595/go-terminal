package drawable

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestBuilder_BasicUnit(t *testing.T) {
	name := "custom-button"

	unit := NewBuilder().
		Name(name).
		Init(func() {}).
		Wipe(func() {}).
		Draw(func(size winsize.Winsize) ([]text.Line, bool) {
			return []text.Line{}, false
		}).
		ToUnit()

	assert.Equal(t, name, unit.Name)
	assert.Len(t, 0, unit.Tags)

	assert.NotNil(t, unit.Drawable.Init)
	assert.NotNil(t, unit.Drawable.Wipe)
	assert.NotNil(t, unit.Drawable.Draw)
}

func TestBuilder_AddTags(t *testing.T) {
	unit := NewBuilder().
		Name("golang").
		AddTags("lang", "google").
		ToUnit()

	assert.GreaterOrEqual(t, 2, len(unit.Tags))
	assert.Contains(t, unit.Tags, "lang")
	assert.Contains(t, unit.Tags, "google")
}

func TestBuilder_MergeTags(t *testing.T) {
	baseTags := set.NewSet[string]()
	baseTags.Add("zig", "c++")

	unit := NewBuilder().
		AddTags("golang").
		MergeTags(baseTags).
		ToUnit()

	assert.GreaterOrEqual(t, 3, len(unit.Tags))
	assert.Contains(t, unit.Tags, "zig")
	assert.Contains(t, unit.Tags, "c++")
	assert.Contains(t, unit.Tags, "golang")
}

func TestBuilder_IncompleteUnit(t *testing.T) {
	unit := NewBuilder().
		Name("empty").
		ToUnit()

	assert.Equal(t, "empty", unit.Name)
	assert.NotNil(t, unit.Tags)

	assert.Nil(t, unit.Drawable.Init)
	assert.Nil(t, unit.Drawable.Wipe)
	assert.Nil(t, unit.Drawable.Draw)
}
