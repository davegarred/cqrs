package cqrs

type EventStore interface {
	Persist(string,[]Event)
	Load(string) []Event
}

type MemEventStore struct {
	eventMap map[string][]Event
}

func (s *MemEventStore) Persist(aggregateId string, newEvents []Event) {
	events := s.eventMap[aggregateId]
	if events == nil{
		events = make([]Event,0)
	}
	events = append(events, newEvents...)
	s.eventMap[aggregateId] = events
}
func (s *MemEventStore) Load(aggregateId string) []Event {
	return s.eventMap[aggregateId]
}

func NewMemEventStore() EventStore {
	return &MemEventStore{make(map[string][]Event)}
}