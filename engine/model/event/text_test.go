package event

import (
	"fmt"
	"testing"
	"time"

	"github.com/Rafael24595/go-terminal/test/support/assert"
	"github.com/Rafael24595/go-terminal/test/support/mock"
)

func applyDeltaStr(buffer string, d *Delta) string {
	return string(ApplyDelta([]rune(buffer), d))
}

func TestForgeEvent_Insert(t *testing.T) {
	s := NewTextEventService()

	m := mergeAction{
		kind:   Insert,
		origin: 5,
		probe:  8,
		extent: 5,
		insert: []string{"a", "b", "c"},
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
		kind:   Insert,
		origin: 5,
		probe:  8,
		extent: 5,
		insert: []string{"a", "b", "c"},
		delete: []string{"A", "Z"},
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
		kind:   DeleteBackward,
		origin: 5,
		probe:  2,
		extent: 5,
		delete: []string{"c", "b", "a"},
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
		kind:   DeleteForward,
		origin: 5,
		probe:  2,
		extent: 5,
		delete: []string{"a", "b", "c"},
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
		kind:   Insert,
		origin: 3,
		probe:  3,
		extent: 7,
		insert: []string{"X"},
		delete: []string{"abcd"},
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
		{kind: Insert, start: 0, end: 0, insert: "g"},
		{kind: Insert, start: 1, end: 1, insert: "o"},
		{kind: Insert, start: 2, end: 2, insert: "l"},
		{kind: Insert, start: 3, end: 3, insert: "a"},
		{kind: Insert, start: 4, end: 4, insert: "n"},
		{kind: Insert, start: 5, end: 5, insert: "g"},
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
		{kind: Insert, start: 0, end: 0, insert: "g"},
		{kind: Insert, start: 2, end: 2, insert: "o"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)

	assert.Equal(t, "g", events[0].insert)
	assert.Equal(t, "o", events[1].insert)
}

func TestMerge_DifferentActions(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: Insert, start: 0, end: 0, insert: "g"},
		{kind: Insert, start: 1, end: 1, insert: "o"},
		{kind: DeleteBackward, start: 1, end: 1, delete: "o"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)
	assert.Equal(t, "go", events[0].insert)
	assert.Equal(t, "o", events[1].delete)
}

func TestMerge_DeleteBackwardContiguous(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: DeleteBackward, start: 5, end: 5, delete: "g"},
		{kind: DeleteBackward, start: 4, end: 4, delete: "i"},
		{kind: DeleteBackward, start: 3, end: 3, delete: "Z"},
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
		{kind: DeleteBackward, start: 5, end: 5},
		{kind: DeleteBackward, start: 2, end: 2},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)
}

func TestMerge_SingleAction(t *testing.T) {
	s := NewTextEventService()

	s.actions = []textAction{
		{kind: Insert, start: 10, end: 10, insert: "Z"},
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

	clock := &mock.TestClock{Time: 1000}
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

	clock := &mock.TestClock{Time: 1000}
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
	clock := &mock.TestClock{Time: 1000}
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
	s.clock = mock.FixedClock(1000)

	s.actions = []textAction{
		{
			kind:      Insert,
			timestamp: 1000 - expires_ms - 1,
		},
	}

	ok := s.shouldFlush(Insert, "a")

	assert.True(t, ok)
}

func TestTextEventService_LimitLogic(t *testing.T) {
	s := NewTextEventService()

	totalPush := event_limit + 50
	for i := range totalPush {
		content := fmt.Sprintf("%d", i)
		s.events = append(s.events, textEvent{
			start:  uint(i),
			delete: content,
			insert: "",
		})
		s.cursor += 1
	}

	s.limitEvents()

	assert.Equal(t, event_limit, len(s.events))
	assert.Equal(t, event_limit, s.cursor)
	assert.Equal(t, "50", s.events[0].delete)

	undoResult := s.Undo()
	assert.NotNil(t, undoResult)
	assert.Equal(t, "249", undoResult.Text)
}

func TestTextEventService_LimitLogicWithPush(t *testing.T) {
	s := NewTextEventService()

	for i := range event_limit + 50 {
		content := fmt.Sprintf("%d", i)
		s.PushEvent(Insert, uint(i), uint(i), content, " ")
	}

	assert.Equal(t, event_limit, len(s.events))
	assert.Equal(t, event_limit, s.cursor)
	assert.Equal(t, "49", s.events[0].delete)

	undoResult := s.Undo()
	assert.NotNil(t, undoResult)
	assert.Equal(t, "249", undoResult.Text)
}

func TestLimitWithCursorAtZero(t *testing.T) {
	s := NewTextEventService()

	for i := range event_limit {
		s.events = append(s.events, textEvent{
			delete: fmt.Sprintf("%d", i)},
		)

		s.cursor++
	}

	for range event_limit {
		s.Undo()
	}

	assert.Equal(t, 0, s.cursor)

	for range 10 {
		s.events = append(s.events, textEvent{
			delete: "new"},
		)

		s.cursor++
	}

	s.limitEvents()

	assert.True(t, s.cursor >= 0)
	assert.Equal(t, 0, s.cursor)
}
