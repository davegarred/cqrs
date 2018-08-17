package cqrs

type Command interface {
	TargetAggregateId() string
}

type Event interface {
	AggregateId() string
}
