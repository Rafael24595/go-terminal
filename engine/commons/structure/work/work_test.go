package work

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestWorkTracker_InitialState(t *testing.T) {
	wt := NewTracker()

	assert.False(t, wt.HasWorks())
	assert.True(t, wt.Finished())
}

func TestWorkTracker_AddTasks(t *testing.T) {
	wt := NewTracker()

	wt.Add(3)

	assert.True(t, wt.HasWorks())
	assert.False(t, wt.Finished())
}

func TestWorkTracker_Advance(t *testing.T) {
	wt := NewTracker()

	wt.Add(2)

	wt.Advance()
	assert.False(t, wt.Finished())

	wt.Advance()
	assert.True(t, wt.Finished())
}

func TestWorkTracker_ResetSafe(t *testing.T) {
	wt := NewTracker()

	works := 5

	wt.Add(works)
	for range works {
		wt.Advance()
	}

	wt.Reset()

	assert.False(t, wt.HasWorks())
	assert.True(t, wt.Finished())
}

func TestWorkTracker_ResetPanic(t *testing.T) {
	wt := NewTracker()

	wt.Add(5)
	wt.Advance()

	assert.Panic(t, func() {
		wt.Reset()
	})
}

func TestWorkTracker_AdvanceWithoutTasks(t *testing.T) {
	wt := NewTracker()

	assert.Panic(t, func() {
		wt.Advance()
	})

	assert.True(t, wt.Finished())
}

func TestWorkTracker_AdvanceOverflow(t *testing.T) {
	wt := NewTracker()

	wt.Add(1)

	wt.Advance()

	assert.Panic(t, func() {
		wt.Advance()
	})

	assert.True(t, wt.Finished())
}

func TestWorkTracker_AddInvalid(t *testing.T) {
	wt := NewTracker()

	assert.Panic(t, func() {
		wt.Add(0)
	})

	assert.Panic(t, func() {
		wt.Add(-5)
	})

	assert.False(t, wt.HasWorks())
	assert.True(t, wt.Finished())
}

func TestWorkTracker_FullFlow(t *testing.T) {
	wt := NewTracker()

	wt.Add(2)
	wt.Add(3)

	for range 5 {
		wt.Advance()
	}

	assert.True(t, wt.Finished())
}
