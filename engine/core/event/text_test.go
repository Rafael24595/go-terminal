package event

import (
	"testing"

	"github.com/Rafael24595/go-terminal/test/support/assert"
)

func TestForgeEvent_Insert(t *testing.T) {
	s := NewTextEventService()

	m := MergeAction{
		action:        Insert,
		initialCaret:  5,
		finalCaret:    8,
		initialAnchor: 5,
		finalAnchor:   8,
		text:          []string{"a", "b", "c"},
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, ev.start, uint(5), "start debe ser 5")
	assert.Equal(t, ev.end, uint(8), "end debe ser 8")
	assert.Equal(t, ev.text, "abc", "texto concatenado")
}

func TestForgeEvent_DeleteBackward(t *testing.T) {
	s := NewTextEventService()

	m := MergeAction{
		action:        DeleteBackward,
		initialCaret:  5,
		finalCaret:    2,
		initialAnchor: 5,
		finalAnchor:   2,
		text:          nil,
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, ev.start, uint(2), "start debe ser 2")
	assert.Equal(t, ev.end, uint(5), "end debe ser 5")
	assert.Equal(t, ev.text, "", "texto vacío en delete")
}

func TestForgeEvent_SelectionActive(t *testing.T) {
	s := NewTextEventService()

	m := MergeAction{
		action:        Insert,
		initialCaret:  3,
		finalCaret:    3,
		initialAnchor: 7,
		finalAnchor:   7,
		text:          []string{"X"},
	}

	ev := s.forgeEvent(m)

	assert.Equal(t, ev.start, uint(3), "start debe ser 3")
	assert.Equal(t, ev.end, uint(7), "end debe ser 7")
	assert.Equal(t, ev.text, "X", "texto insertado")
}

func TestMergeActions_MultipleInserts(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 0, anchor: 0, text: "g"},
		{action: Insert, caret: 1, anchor: 1, text: "o"},
		{action: Insert, caret: 2, anchor: 2, text: "l"},
		{action: Insert, caret: 3, anchor: 3, text: "a"},
		{action: Insert, caret: 4, anchor: 4, text: "n"},
		{action: Insert, caret: 5, anchor: 5, text: "g"},
	}

	events := s.mergeActions(s.actions)

	assert.Equal(t, len(events), 1)

	ev := events[0]

	assert.Equal(t, ev.start, uint(0))
	assert.Equal(t, ev.end, uint(5))
	assert.Equal(t, ev.text, "golang")
}

func TestMerge_InsertNonContiguous(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 0, anchor: 0, text: "g"},
		{action: Insert, caret: 2, anchor: 2, text: "o"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)

	assert.Equal(t, "g", events[0].text)
	assert.Equal(t, "o", events[1].text)
}

func TestMerge_DifferentActions(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 0, anchor: 0, text: "g"},
		{action: DeleteBackward, caret: 1, anchor: 1},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)
	assert.Equal(t, Insert, events[0].action)
	assert.Equal(t, DeleteBackward, events[1].action)
}

func TestMerge_DeleteBackwardContiguous(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: DeleteBackward, caret: 5, anchor: 5},
		{action: DeleteBackward, caret: 4, anchor: 4},
		{action: DeleteBackward, caret: 3, anchor: 3},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 1, events)

	ev := events[0]
	assert.Equal(t, uint(3), ev.start)
	assert.Equal(t, uint(5), ev.end)
}

func TestMerge_DeleteBackwardNonContiguous(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: DeleteBackward, caret: 5, anchor: 5},
		{action: DeleteBackward, caret: 2, anchor: 2},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 2, events)
}

// func TestMerge_AnchorChangeBreaksMerge(t *testing.T) {
// 	s := NewTextEventService()

// 	s.actions = []TextAction{
// 		{action: Insert, caret: 0, anchor: 0, text: "a"},
// 		{action: Insert, caret: 1, anchor: 2, text: "b"},
// 	}

// 	events := s.mergeActions(s.actions)

// 	assert.Len(t, 2, events)
// }

func TestMerge_SingleAction(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 10, anchor: 10, text: "Z"},
	}

	events := s.mergeActions(s.actions)

	assert.Len(t, 1, events)
	assert.Equal(t, "Z", events[0].text)
}

func TestApplyLastEvent_NoActionsNoEvents(t *testing.T) {
	s := NewTextEventService()

	buff := []rune("hello")
	out := s.ApplyLastEvent(buff)

	assert.Equal(t, "hello", string(out))
}

func TestApplyLastEvent_Insert(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 5, anchor: 5, text: " world"},
	}

	buff := []rune("hello")
	out := s.ApplyLastEvent(buff)

	assert.Equal(t, "hello world", string(out))
	assert.Len(t, 0, s.actions)
	assert.Len(t, 0, s.events)
}

func TestApplyLastEvent_MergedInsert(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 0, anchor: 0, text: "a"},
		{action: Insert, caret: 1, anchor: 1, text: "b"},
		{action: Insert, caret: 2, anchor: 2, text: "c"},
	}

	buff := make([]rune, 2)
	out := s.ApplyLastEvent(buff)

	assert.Equal(t, "abc", string(out))
	assert.Len(t, 0, s.events)
}

func TestApplyLastEvent_DeleteBackward(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: DeleteBackward, caret: 4, anchor: 5},
	}

	buff := []rune("hello")
	out := s.ApplyLastEvent(buff)

	assert.Equal(t, "hell", string(out))
}

func TestApplyLastEvent_ReplaceSelection(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 1, anchor: 4, text: "i"},
	}

	buff := []rune("hello")
	out := s.ApplyLastEvent(buff)

	assert.Equal(t, "hio", string(out))
}

func TestApplyLastEvent_LIFO(t *testing.T) {
	s := NewTextEventService()

	s.events = []TextEvent{
		{action: Insert, start: 2, end: 2, text: "A"},
		{action: Insert, start: 1, end: 1, text: "B"},
	}

	buff := []rune("x")
	out := s.ApplyLastEvent(buff)

	assert.Equal(t, "xB", string(out))
	assert.Len(t, 1, s.events)

	out = s.ApplyLastEvent(out)

	assert.Equal(t, "xBA", string(out))
	assert.Len(t, 0, s.events)
}

func TestApplyLastEvent_ActionsConsumed(t *testing.T) {
	s := NewTextEventService()

	s.actions = []TextAction{
		{action: Insert, caret: 0, anchor: 0, text: "a"},
	}

	buff := []rune("")
	out1 := s.ApplyLastEvent(buff)
	out2 := s.ApplyLastEvent(out1)

	assert.Equal(t, "a", string(out1))
	assert.Equal(t, "a", string(out2))
}
