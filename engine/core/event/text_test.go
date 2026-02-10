package event

import (
	"testing"
	"time"

	"github.com/Rafael24595/go-terminal/test/support/assert"
)

type TestClock struct {
	now int64
}

func (c *TestClock) Now() int64 {
	return c.now
}

func (c *TestClock) Advance(ms int64) {
	c.now += ms
}

func fixedClock(t int64) clock {
	return func() int64 {
		return t
	}
}

func applyDeltaStr(buffer string, d *Delta) string {
	return string(ApplyDelta([]rune(buffer), d))
}

func TestForgeEvent_Insert(t *testing.T) {
	s := NewTextEventService()

	m := mergeAction{
		kind:          Insert,
		initialCaret:  5,
		finalCaret:    8,
		initialAnchor: 5,
		finalAnchor:   8,
		insert:        []string{"a", "b", "c"},
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, uint(5), ev.start)
	assert.Equal(t, uint(8), ev.start+uint(len(ev.insert)))
	assert.Equal(t, "abc", ev.insert)
	assert.Equal(t, "", ev.delete)
}

func TestForgeEvent_Replace(t *testing.T) {
	s := NewTextEventService()

	m := mergeAction{
		kind:          Insert,
		initialCaret:  5,
		finalCaret:    8,
		initialAnchor: 5,
		finalAnchor:   8,
		insert:        []string{"a", "b", "c"},
		delete:        []string{"A", "Z"},
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, uint(5), ev.start)
	assert.Equal(t, uint(8), ev.start+uint(len(ev.insert)))
	assert.Equal(t, "abc", ev.insert)
	assert.Equal(t, "AZ", ev.delete)
}

func TestForgeEvent_DeleteBackward(t *testing.T) {
	s := NewTextEventService()

	m := mergeAction{
		kind:          DeleteBackward,
		initialCaret:  5,
		finalCaret:    2,
		initialAnchor: 5,
		finalAnchor:   2,
		delete:        []string{"c", "b", "a"},
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, uint(2), ev.start)
	assert.Equal(t, uint(5), ev.start+uint(len(ev.delete)))
	assert.Equal(t, "abc", ev.delete)
	assert.Equal(t, "", ev.insert)
}

func TestForgeEvent_DeleteForward(t *testing.T) {
	s := NewTextEventService()

	m := mergeAction{
		kind:          DeleteForward,
		initialCaret:  5,
		finalCaret:    2,
		initialAnchor: 5,
		finalAnchor:   2,
		delete:        []string{"a", "b", "c"},
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, uint(2), ev.start)
	assert.Equal(t, uint(5), ev.start+uint(len(ev.delete)))
	assert.Equal(t, "abc", ev.delete)
	assert.Equal(t, "", ev.insert)
}

func TestForgeEvent_SelectionActive(t *testing.T) {
	s := NewTextEventService()

	m := mergeAction{
		kind:          Insert,
		initialCaret:  3,
		finalCaret:    3,
		initialAnchor: 7,
		finalAnchor:   7,
		insert:        []string{"X"},
		delete:        []string{"abcd"},
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, uint(3), ev.start)
	assert.Equal(t, uint(7), ev.start+uint(len(ev.delete)))

	assert.Equal(t, "X", ev.insert)
	assert.Equal(t, "abcd", ev.delete)
}

func TestMergeActions_MultipleInserts(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: Insert, caret: 0, anchor: 0, insert: "g"},
		{kind: Insert, caret: 1, anchor: 1, insert: "o"},
		{kind: Insert, caret: 2, anchor: 2, insert: "l"},
		{kind: Insert, caret: 3, anchor: 3, insert: "a"},
		{kind: Insert, caret: 4, anchor: 4, insert: "n"},
		{kind: Insert, caret: 5, anchor: 5, insert: "g"},
	}

	events := s.mergeActions(s.actions)

	assert.Equal(t, len(events), 1)

	ev := events[0]

	assert.Equal(t, uint(0), ev.start)
	assert.Equal(t, uint(6), ev.start+uint(len(ev.insert)))
	assert.Equal(t, "golang", ev.insert)
}

func TestMerge_InsertNonContiguous(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: Insert, caret: 0, anchor: 0, insert: "g"},
		{kind: Insert, caret: 2, anchor: 2, insert: "o"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)

	assert.Equal(t, "g", events[0].insert)
	assert.Equal(t, "o", events[1].insert)
}

func TestMerge_DifferentActions(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: Insert, caret: 0, anchor: 0, insert: "g"},
		{kind: Insert, caret: 1, anchor: 1, insert: "o"},
		{kind: DeleteBackward, caret: 1, anchor: 1, delete: "o"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)
	assert.Equal(t, "go", events[0].insert)
	assert.Equal(t, "o", events[1].delete)
}

func TestMerge_DeleteBackwardContiguous(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: DeleteBackward, caret: 5, anchor: 5, delete: "g"},
		{kind: DeleteBackward, caret: 4, anchor: 4, delete: "i"},
		{kind: DeleteBackward, caret: 3, anchor: 3, delete: "Z"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 1, events)

	ev := events[0]
	assert.Equal(t, uint(3), ev.start)
	assert.Equal(t, uint(6), ev.start+uint(len(ev.delete)))
	assert.Equal(t, "Zig", ev.delete)
}

func TestMerge_DeleteBackwardNonContiguous(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: DeleteBackward, caret: 5, anchor: 5},
		{kind: DeleteBackward, caret: 2, anchor: 2},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)
}

func TestMerge_SingleAction(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: Insert, caret: 10, anchor: 10, insert: "Z"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 1, events)
	assert.Equal(t, "Z", events[0].insert)
}

func TestShouldFlush_NoActions(t *testing.T) {
	s := NewTextEventService()

	ok := s.shouldFlush(Insert, "a")

	assert.False(t, ok)
}

func TestShouldFlush_SameAction_NoSpace_NotExpired(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{
			kind:      Insert,
			insert:    "a",
			timestamp: time.Now().UnixMilli(),
		},
	}

	ok := s.shouldFlush(Insert, "b")

	assert.False(t, ok)
}

func TestShouldFlush_ActionChange(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{
			kind:      Insert,
			timestamp: time.Now().UnixMilli(),
		},
	}

	ok := s.shouldFlush(DeleteBackward, "")

	assert.True(t, ok)
}

func TestShouldFlush_Whitespace(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{
			kind:      Insert,
			timestamp: time.Now().UnixMilli(),
		},
	}

	ok := s.shouldFlush(Insert, " ")

	assert.True(t, ok)
}

func TestShouldFlush_Newline(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{
			kind:      Insert,
			timestamp: time.Now().UnixMilli(),
		},
	}

	ok := s.shouldFlush(Insert, "\n")

	assert.True(t, ok)
}

func TestShouldFlush_Expired(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{
			kind:      Insert,
			timestamp: time.Now().UnixMilli() - expires_ms - 1,
		},
	}

	ok := s.shouldFlush(Insert, "a")

	assert.True(t, ok)
}

func TestPushEvent_AddsAction(t *testing.T) {
	s := NewTextEventService()

	s.PushEvent(Insert, 0, 0, "", "a")

	assert.Len(t, 1, s.actions)
	assert.Equal(t, Insert, s.actions[0].kind)
	assert.Equal(t, "a", s.actions[0].insert)
}

func TestPushEvent_FlushOnWhitespace(t *testing.T) {
	s := NewTextEventService()

	s.PushEvent(Insert, 0, 0, "", "a")
	s.PushEvent(Insert, 1, 1, "", " ")

	assert.Len(t, 1, s.actions)
	assert.Len(t, 1, s.events)
}

func TestPushEvent_FlushOnActionChange(t *testing.T) {
	s := NewTextEventService()

	s.PushEvent(Insert, 0, 0, "", "a")
	s.PushEvent(DeleteBackward, 1, 1, "a", "")

	assert.Len(t, 1, s.actions)
	assert.Len(t, 1, s.events)
}

func TestPushEvent_FlushOnExpire(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{
			kind:      Insert,
			insert:    "a",
			timestamp: time.Now().UnixMilli() - expires_ms - 1,
		},
	}

	s.PushEvent(Insert, 1, 1, "", "b")

	assert.Len(t, 1, s.actions)
	assert.Len(t, 1, s.events)
}

func TestPushEvent_Typing(t *testing.T) {
	s := NewTextEventService()

	clock := &TestClock{now: 1000}
	s.clock = clock.Now

	i := 0
	for _, v := range "Golang" {
		s.PushEvent(Insert, uint(i), uint(i), "", string(v))
		clock.Advance(100)
		i++
	}

	s.PushEvent(Insert, uint(i), uint(i), "", " ")
	i++

	for _, v := range "Zig" {
		s.PushEvent(Insert, uint(i), uint(i), "", string(v))
		clock.Advance(expires_ms + 1)
		i++
	}

	s.PushEvent(Insert, uint(i), uint(i), "", " ")
	i++

	assert.Len(t, 1, s.actions)
	assert.Len(t, 4, s.events)

	assert.Equal(t, s.events[0].insert, "Golang")
	assert.Equal(t, s.events[1].insert, " "+"Z")
	assert.Equal(t, s.events[2].insert, "i")
	assert.Equal(t, s.events[3].insert, "g")
}

func TestPushEvent_UndoAndRedo(t *testing.T) {
	s := NewTextEventService()

	clock := &TestClock{now: 1000}
	s.clock = clock.Now

	i := 0
	for _, v := range "Golang" {
		s.PushEvent(Insert, uint(i), uint(i), "", string(v))
		clock.Advance(100)
		i++
	}

	s.PushEvent(Insert, uint(i), uint(i), "", " ")
	i++

	clock.Advance(expires_ms + 1)

	for _, v := range "Zig" {
		s.PushEvent(Insert, uint(i), uint(i), "", string(v))
		clock.Advance(100)
		i++
	}

	assert.Len(t, 3, s.actions)
	assert.Len(t, 2, s.events)

	buff := "Golang Zig"

	evnt := s.Undo()
	assert.NotNil(t, evnt)

	buff = applyDeltaStr(buff, evnt)
	assert.Equal(t, "Golang ", buff)

	evnt = s.Redo()
	assert.NotNil(t, evnt)

	buff = applyDeltaStr(buff, evnt)
	assert.Equal(t, "Golang Zig", buff)
}

func TestPushEvent_UndoRedoTruncateHistory(t *testing.T) {
	s := NewTextEventService()
	clock := &TestClock{now: 1000}
	s.clock = clock.Now

	i := 0
	for _, v := range "Golang " {
		s.PushEvent(Insert, uint(i), uint(i), "", string(v))
		clock.Advance(100)
		i++
	}

	clock.Advance(expires_ms + 1)
	for _, v := range "Zig" {
		s.PushEvent(Insert, uint(i), uint(i), "", string(v))
		clock.Advance(100)
		i++
	}

	buff := "Golang Zig"

	evnt := s.Undo()
	assert.NotNil(t, evnt)

	buff = applyDeltaStr(buff, evnt)
	assert.Equal(t, "Golang ", string(buff))
	i = len(buff)

	s.PushEvent(Insert, uint(i), uint(i), "", "New")
	assert.Equal(t, s.cursor, len(s.events))

	_ = s.Undo()

	evnt = s.Redo()
	assert.NotNil(t, evnt)

	buff = applyDeltaStr(buff, evnt)
	assert.Equal(t, "Golang New", string(buff))
}

func TestPushEvent_UndoRedoHistoryConsistence(t *testing.T) {
	s := NewTextEventService()

	s.PushEvent(DeleteForward, uint(7), uint(11), "Rust ", "")

	evnt := s.Undo()
	assert.NotNil(t, evnt)

	buff := "Golang Zig"
	buff = applyDeltaStr(buff, evnt)
	assert.Equal(t, "Golang Rust Zig", string(buff))

	evnt = s.Redo()
	assert.NotNil(t, evnt)

	buff = applyDeltaStr(buff, evnt)
	assert.Equal(t, "Golang Zig", string(buff))
}

func TestPushEvent_UndoRedoHistoryConsistenceWithLoop(t *testing.T) {
	s := NewTextEventService()

	buff := "Golang Rust Zig"

	s.PushEvent(Insert, 7, 12, "X ", "Rust ")

	for range 10 {
		buff = applyDeltaStr(buff, s.Undo())
		assert.Equal(t, "Golang X Zig", buff)

		buff = applyDeltaStr(buff, s.Redo())
		assert.Equal(t, "Golang Rust Zig", buff)
	}
}

func TestShouldFlush_Expired_WithClock(t *testing.T) {
	s := NewTextEventService()
	s.clock = fixedClock(1000)

	s.actions = []textAction{
		{
			kind:      Insert,
			timestamp: 1000 - expires_ms - 1,
		},
	}

	ok := s.shouldFlush(Insert, "a")

	assert.True(t, ok)
}
