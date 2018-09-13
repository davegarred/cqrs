package persist

import (
	"encoding/json"
	"github.com/davegarred/cqrs"
	"reflect"
)

type MemEventStore struct {
	eventBus cqrs.EventBus
	eventMap map[string][]StoredEvent
}

type StoredEvent struct {
	eventType reflect.Type
	payload   []byte
}

func (s *MemEventStore) Persist(aggregateId string, newEvents []cqrs.Event) {
	events := s.eventMap[aggregateId]
	if events == nil {
		events = make([]StoredEvent, 0)
	}
	for _, event := range newEvents {
		events = append(events, serialize(event))
	}
	s.eventMap[aggregateId] = events
	s.eventBus.PublishEvents(newEvents)
}

func (s *MemEventStore) Load(aggregateId string) []cqrs.Event {
	storedEvents := s.eventMap[aggregateId]
	events := make([]cqrs.Event, len(storedEvents))
	for i, storedEvent := range storedEvents {
		event := reflect.New(storedEvent.eventType).Interface().(cqrs.Event)
		err := json.Unmarshal(storedEvent.payload, event)
		if err != nil {
			panic(err)
		}
		events[i] = reflect.ValueOf(event).Elem().Interface().(cqrs.Event)
	}
	return events
}

func NewMemEventStore(eventBus cqrs.EventBus) cqrs.EventStore {
	return &MemEventStore{eventBus, make(map[string][]StoredEvent)}
}

func serialize(event cqrs.Event) StoredEvent {
	eventType := reflect.TypeOf(event)
	payload, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	return StoredEvent{eventType, payload}
}
