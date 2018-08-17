package cqrs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const aggregateId = "aggregate_id"

func TestNewMemEventStore(t *testing.T) {
	es := NewMemEventStore()
	fooCreatedEvent := FooCreatedEvent{aggregateId}
	fooNamedEvent := FooNamedEvent{aggregateId, "aname forfoo"}
	es.Persist(aggregateId, []Event{fooCreatedEvent, fooNamedEvent})

	events := es.Load(aggregateId)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, fooCreatedEvent, events[0])
	assert.Equal(t, fooNamedEvent, events[1])
}
