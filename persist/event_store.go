package persist

import (
	"encoding/json"
	"github.com/davegarred/cqrs/ext"
	"reflect"
)

type MemEventStore struct {
	eventBus ext.EventBus
	eventMap map[string][]StoredEvent
}

type StoredEvent struct {
	eventType reflect.Type
	payload   []byte
}

func (s *MemEventStore) Persist(aggregateId string, newEvents []ext.Event) {
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

func (s *MemEventStore) Load(aggregateId string) []ext.Event {
	storedEvents := s.eventMap[aggregateId]
	events := make([]ext.Event, len(storedEvents))
	for i, storedEvent := range storedEvents {
		event := reflect.New(storedEvent.eventType).Interface().(ext.Event)
		err := json.Unmarshal(storedEvent.payload, event)
		if err != nil {
			panic(err)
		}
		events[i] = reflect.ValueOf(event).Elem().Interface().(ext.Event)
	}
	return events
}

func NewMemEventStore(eventBus ext.EventBus) ext.EventStore {
	return &MemEventStore{eventBus, make(map[string][]StoredEvent)}
}

func serialize(event ext.Event) StoredEvent {
	eventType := reflect.TypeOf(event)
	payload, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	return StoredEvent{eventType, payload}
}
