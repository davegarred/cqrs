package cqrs

import (
	"encoding/json"
	"reflect"
)

type EventStore interface {
	Persist(aggregateId string, events []Event)
	Load(aggregateId string) []Event
}

type MemEventStore struct {
	eventMap map[string][]StoredEvent
}

type StoredEvent struct {
	eventType reflect.Type
	payload   []byte
}

func (s *MemEventStore) Persist(aggregateId string, newEvents []Event) {
	events := s.eventMap[aggregateId]
	if events == nil {
		events = make([]StoredEvent, 0)
	}
	for _, event := range newEvents {
		events = append(events, serialize(event))
	}
	s.eventMap[aggregateId] = events
}
func (s *MemEventStore) Load(aggregateId string) []Event {
	storedEvents := s.eventMap[aggregateId]
	events := make([]Event, len(storedEvents))
	for i, storedEvent := range storedEvents {
		event := reflect.New(storedEvent.eventType).Interface().(Event)
		err := json.Unmarshal(storedEvent.payload, event)
		if err != nil {
			panic(err)
		}
		events[i] = reflect.ValueOf(event).Elem().Interface().(Event)
	}
	return events
}

func NewMemEventStore() EventStore {
	return &MemEventStore{make(map[string][]StoredEvent)}
}

func serialize(event Event) StoredEvent {
	eventType := reflect.TypeOf(event)
	payload, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	return StoredEvent{eventType, payload}
}
