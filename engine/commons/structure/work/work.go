package work

import assert "github.com/Rafael24595/go-assert/assert/runtime"

type Tracker struct {
	total  uint
	cursor uint
}

func NewTracker() *Tracker {
	return &Tracker{}
}

func (t *Tracker) Add(tasks int) *Tracker {
	if tasks <= 0 {
		assert.Unreachable("tasks should be greater than 0")

		return t
	}

	t.total += uint(tasks)
	return t
}

func (t *Tracker) Advance() *Tracker {
	if t.cursor >= t.total {
		assert.Unreachable("task cursor overflow %d/%d", t.cursor, t.total)

		t.cursor = t.total
		return t
	}

	t.cursor++
	return t
}

func (t *Tracker) Reset() *Tracker {
	assert.True(t.cursor == t.total, "invalid state %d/%d", t.cursor, t.total)

	t.total = 0
	t.cursor = 0
	return t
}

func (t *Tracker) HasWorks() bool {
	return t.total > 0
}

func (t *Tracker) Finished() bool {
	return t.cursor >= t.total
}
