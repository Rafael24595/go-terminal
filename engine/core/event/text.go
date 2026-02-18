package event

import (
	"strings"
	"time"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
)

const expires_ms = 1000

const event_limit = 200
const action_limit = 2000

type ActionKind int

const (
	Insert ActionKind = iota
	DeleteBackward
	DeleteForward
)

type textAction struct {
	kind      ActionKind
	start     uint
	end       uint
	delete    string
	insert    string
	timestamp int64
}

type mergeAction struct {
	kind   ActionKind
	origin uint
	extent uint
	probe  uint
	delete []string
	insert []string
}

func (m *mergeAction) len() uint {
	var n uint
	for _, t := range m.insert {
		r := []rune(t)
		n += uint(len(r))
	}
	return n
}

type textEvent struct {
	start  uint
	insert string
	delete string
}

type Delta struct {
	Start uint
	End   uint
	Text  string
}

type clock func() int64

func unixClock() int64 {
	return time.Now().UnixMilli()
}

type TextEventService struct {
	clock   clock
	actions []textAction
	events  []textEvent
	cursor  int
}

func NewTextEventService() *TextEventService {
	return &TextEventService{
		clock:   unixClock,
		actions: make([]textAction, 0),
		events:  make([]textEvent, 0, event_limit),
	}
}

func (s *TextEventService) mergeActions(actions []textAction) []textEvent {
	events := make([]textEvent, 0)

	var event *mergeAction = nil

	i := 0
	for i < len(actions) {
		action := actions[i]

		if event == nil {
			event = &mergeAction{
				kind:   action.kind,
				origin: action.start,
				extent: action.end,
				probe:  action.start,
				delete: []string{action.delete},
				insert: []string{action.insert},
			}

			i++
			continue
		}

		isConsistent := s.isConsistentAction(action, *event)
		if action.kind != event.kind || !isConsistent {
			events = append(events, s.forgeEvent(*event))
			event = nil

			continue
		}

		event.delete = append(event.delete, action.delete)
		event.insert = append(event.insert, action.insert)

		event.probe = action.start

		i++
	}

	if event != nil {
		events = append(events, s.forgeEvent(*event))
	}

	return events
}

func (s *TextEventService) isConsistentAction(action textAction, event mergeAction) bool {
	switch action.kind {
	case Insert:
		return event.origin+event.len() == action.start
	case DeleteBackward:
		return action.start == math.SubClampZero(event.probe, uint(1))
	case DeleteForward:
		return action.start == event.probe+1
	}

	return false
}

func (s *TextEventService) forgeEvent(action mergeAction) textEvent {
	start := min(
		action.origin,
		action.probe,
		action.extent,
	)

	insert := ""
	delete := ""

	switch action.kind {
	case Insert, DeleteForward:
		insert = strings.Join(action.insert, "")
		delete = strings.Join(action.delete, "")

	case DeleteBackward:
		insert = runes.JoinReverse(action.insert)
		delete = runes.JoinReverse(action.delete)
	}

	return textEvent{
		start:  start,
		insert: insert,
		delete: delete,
	}
}

func (s *TextEventService) PushEvent(action ActionKind, start uint, end uint, delete, insert string) {
	s.events = s.events[:s.cursor]

	if s.shouldFlush(action, insert) {
		s.flushAndLimit()
	}

	now := s.clock()

	s.actions = append(s.actions, textAction{
		kind:      action,
		start:     start,
		end:       end,
		delete:    delete,
		insert:    insert,
		timestamp: now,
	})
}

func (s *TextEventService) Undo() *Delta {
	s.flushAndLimit()

	if len(s.events) == 0 || s.cursor == 0 {
		return nil
	}

	s.decrementCursor()

	event := s.events[s.cursor]

	return &Delta{
		Start: event.start,
		End:   event.start + uint(len([]rune(event.insert))),
		Text:  event.delete,
	}
}

func (s *TextEventService) Redo() *Delta {
	s.flushAndLimit()

	if len(s.events) == 0 || s.cursor >= len(s.events) {
		return nil
	}

	event := s.events[s.cursor]

	s.incrementCursor()

	return &Delta{
		Start: event.start,
		End:   event.start + uint(len([]rune(event.delete))),
		Text:  event.insert,
	}
}

func (s *TextEventService) incrementCursor() {
	s.cursor = min(len(s.events), s.cursor+1)
}

func (s *TextEventService) decrementCursor() {
	s.cursor = max(0, s.cursor-1)
}

func (s *TextEventService) flushAndLimit() {
	s.flushActions()
	s.limitEvents()
}

func (s *TextEventService) flushActions() {
	if len(s.actions) == 0 {
		return
	}

	events := s.mergeActions(s.actions)
	s.events = append(s.events, events...)
	s.cursor = len(s.events)
	s.actions = nil
}

func (s *TextEventService) shouldFlush(action ActionKind, text string) bool {
	if len(s.actions) > action_limit {
		return false
	}

	if len(s.actions) == 0 {
		return false
	}

	if strings.ContainsAny(text, " \n") {
		return true
	}

	last := s.actions[len(s.actions)-1]
	if last.kind != action {
		return true
	}

	now := s.clock()
	if now-last.timestamp >= expires_ms {
		return true
	}

	return false
}

func (s *TextEventService) limitEvents() {
	if len(s.events) <= event_limit {
		return
	}

	buff := make([]textEvent, event_limit)

	excess := len(s.events) - event_limit
	copy(buff, s.events[excess:])

	s.events = buff
	s.cursor = max(0, s.cursor-excess)
}

func ApplyDelta(insert []rune, d *Delta) []rune {
	size := uint(len(insert))
	if d.Start > size || d.End > size {
		return insert
	}

	runes := []rune(d.Text)
	runesSize := uint(len(runes))

	tail := size - d.End
	total := d.Start + uint(len(runes)) + tail

	newBuffer := make([]rune, total)

	copy(newBuffer[:d.Start], insert[:d.Start])
	copy(newBuffer[d.Start:], runes)
	copy(newBuffer[d.Start+runesSize:], insert[d.End:])

	return newBuffer
}
