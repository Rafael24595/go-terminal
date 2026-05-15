package pipeline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestExpiration_Persistent_NeverExpires(t *testing.T) {
	exp := persistent()

	mock1 := screen_test.MockScreen{
		Name: "A",
	}.ToNode()

	mock2 := screen_test.MockScreen{
		Name: "B",
	}.ToNode()

	assert.False(t, exp.on(&mock1))
	assert.False(t, exp.on(&mock2))
}

func TestExpiration_OnNode_SameNode_DoesNotExpire(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "A",
	}.ToNode()

	exp := onNode(&mock)

	assert.False(t, exp.on(&mock))
}

func TestExpiration_OnNode_DifferentNode_Expires(t *testing.T) {
	mock1 := screen_test.MockScreen{
		Name: "A",
	}.ToNode()

	mock2 := screen_test.MockScreen{
		Name: "B",
	}.ToNode()

	exp := onNode(&mock1)

	assert.True(t, exp.on(&mock2))
}

func TestExpiration_OnName_SameName_Expires(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "A",
	}.ToNode()

	exp := onName("A")

	assert.True(t, exp.on(&mock))
}

func TestExpiration_OnName_DifferentName_DoesNotExpire(t *testing.T) {
	mock := screen_test.MockScreen{
		Name: "A",
	}.ToNode()

	exp := onName("B")

	assert.False(t, exp.on(&mock))
}
