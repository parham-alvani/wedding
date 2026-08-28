package model

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Event identifies one of the ceremonies a guest can be invited to. The site
// has a page per event, so the set is closed rather than user-defined.
type Event string

const (
	EventEngagement Event = "engagement"
	EventWedding    Event = "wedding"
)

// eventSeparator joins events when they are stored or written to a CSV.
const eventSeparator = ","

var ErrUnknownEvent = errors.New("unknown event")

// AllEvents lists every event, in the order they happen.
func AllEvents() []Event {
	return []Event{EventEngagement, EventWedding}
}

// ParseEvent validates a single event name.
func ParseEvent(raw string) (Event, error) {
	event := Event(strings.ToLower(strings.TrimSpace(raw)))

	if slices.Contains(AllEvents(), event) {
		return event, nil
	}

	return "", fmt.Errorf("%w %q: expected one of %s",
		ErrUnknownEvent, raw, strings.Join(EventNames(AllEvents()), ", "))
}

// ParseEvents reads a separated list such as "engagement,wedding". An empty
// string means every event, which is what an invitation with no tiers is.
func ParseEvents(raw string) ([]Event, error) {
	if strings.TrimSpace(raw) == "" {
		return AllEvents(), nil
	}

	seen := make(map[Event]bool, len(AllEvents()))

	for _, part := range strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '|' || r == ' '
	}) {
		event, err := ParseEvent(part)
		if err != nil {
			return nil, err
		}

		seen[event] = true
	}

	if len(seen) == 0 {
		return AllEvents(), nil
	}

	return sortEvents(seen), nil
}

// FormatEvents renders events for storage, always in ceremony order so the
// stored value does not depend on how it was typed.
func FormatEvents(events []Event) string {
	seen := make(map[Event]bool, len(events))
	for _, event := range events {
		seen[event] = true
	}

	return strings.Join(EventNames(sortEvents(seen)), eventSeparator)
}

// EventNames converts events to plain strings.
func EventNames(events []Event) []string {
	out := make([]string, 0, len(events))
	for _, event := range events {
		out = append(out, string(event))
	}

	return out
}

func sortEvents(seen map[Event]bool) []Event {
	out := make([]Event, 0, len(seen))

	for _, known := range AllEvents() {
		if seen[known] {
			out = append(out, known)
		}
	}

	return out
}
