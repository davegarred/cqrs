package cqrs

type Command interface {
	TargetAggregateId() string
}

type Event interface {
	AggregateId() string
}

type EventStore interface {
	Persist(aggregateId string, events []Event)
	Load(aggregateId string) []Event
}
