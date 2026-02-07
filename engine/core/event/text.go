package event

import (
	"strings"

	"github.com/Rafael24595/go-terminal/engine/helper/math"
	"github.com/Rafael24595/go-terminal/engine/helper/runes"
)

type ActionKind int

const (
	Insert ActionKind = iota
	DeleteBackward
	DeleteForward
)

type TextAction struct {
	action    ActionKind
	caret     uint
	anchor    uint
	text      string
	timestamp int64
}

type MergeAction struct {
	action        ActionKind
	initialCaret  uint
	initialAnchor uint
	finalCaret    uint
	finalAnchor   uint
	text          []string
	timestamp     int64
}

func (m *MergeAction) len() uint {
	var n uint
	for _, t := range m.text {
		n += uint(len(t))
	}
	return n
}

type TextEvent struct {
	action ActionKind
	start  uint
	end    uint
	text   string
}

type TextEventService struct {
	actions []TextAction
	events  []TextEvent
}

func NewTextEventService() *TextEventService {
	return &TextEventService{
		actions: make([]TextAction, 0),
		events:  make([]TextEvent, 0),
	}
}

func (s *TextEventService) mergeActions(actions []TextAction) []TextEvent {
	events := make([]TextEvent, 0)

	var event *MergeAction = nil

	i := 0
	for i < len(actions) {
		action := actions[i]

		if event == nil {
			event = &MergeAction{
				action:        action.action,
				initialCaret:  action.caret,
				initialAnchor: action.anchor,
				finalCaret:    action.caret,
				finalAnchor:   action.anchor,
				text:          []string{action.text},
			}

			i++
			continue
		}

		isConsistent := s.isConsistentAction(action, *event)
		if action.action != event.action || !isConsistent {
			events = append(events, s.forgeEvent(*event))
			event = nil

			continue
		}

		event.text = append(event.text, action.text)
		event.finalCaret = action.caret
		event.finalAnchor = action.anchor

		i++
	}

	if event != nil {
		events = append(events, s.forgeEvent(*event))
	}

	return events
}

func (s *TextEventService) isConsistentAction(action TextAction, event MergeAction) bool {
	switch action.action {
	case Insert:
		return event.initialCaret+event.len() == action.caret
	case DeleteBackward:
		return action.caret == math.SubClampZero(event.finalCaret, uint(1))
	case DeleteForward:
		return action.caret == event.finalCaret+1
	}

	return false
}

func (s *TextEventService) forgeEvent(action MergeAction) TextEvent {
	start := min(
		action.initialCaret,
		action.finalCaret,
		action.initialAnchor,
		action.finalAnchor,
	)

	end := max(
		action.initialCaret,
		action.finalCaret,
		action.initialAnchor,
		action.finalAnchor,
	)

	return TextEvent{
		action: action.action,
		start:  start,
		end:    end,
		text:   s.joinText(action),
	}
}

func (s *TextEventService) joinText(action MergeAction) string {
	if action.action == Insert {
		return strings.Join(action.text, "")
	}

	return ""
}

func (s *TextEventService) PushEvent(action TextAction, caret uint, anchor uint, text string) {
	//
}

func (s *TextEventService) ApplyLastEvent(buff []rune) []rune {
	s.FlushActions()

	if len(s.events) == 0 {
		return buff
	}

	event := s.events[len(s.events)-1]
	s.events = s.events[:len(s.events)-1]

	return apply(buff, event)
}

func (s *TextEventService) FlushActions() {
	if len(s.actions) == 0 {
		return
	}

	events := s.mergeActions(s.actions)
	s.events = append(s.events, events...)
	s.actions = nil
}

func apply(buff []rune, ev TextEvent) []rune {
	return runes.AppendRange(
		buff,
		[]rune(ev.text),
		ev.start,
		ev.end,
	)
}
