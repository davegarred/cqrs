package persist

import (
	"github.com/davegarred/cqrs"
	"github.com/davegarred/cqrs/ext"
	"testing"

	"github.com/stretchr/testify/assert"
)

const aggregateId = "aggregate_id"

func TestNewMemEventStore(t *testing.T) {
	assert := assert.New(t)
	listener := &eventBusQueryListener{}
	eventBus := cqrs.NewEventBus()
	eventBus.RegisterQueryEventHandlers(listener)
	es := NewMemEventStore(eventBus)
	event1 := eventBusTestEvent1{aggregateId}
	event2 := eventBusTestEvent2{aggregateId, "a name"}

	es.Persist(aggregateId, []ext.Event{event1, event2})

	events := es.Load(aggregateId)
	assert.Equal(2, len(events))
	assert.Equal(event1, events[0])
	assert.Equal(event2, events[1])
	assert.True(listener.foundEvent1)
	assert.True(listener.foundEvent2)
}

type eventBusTestEvent1 struct {
	Id string
}

func (e eventBusTestEvent1) AggregateId() string { return e.Id }

type eventBusTestEvent2 struct {
	Id   string
	Name string
}

func (e eventBusTestEvent2) AggregateId() string { return e.Id }

type eventBusQueryListener struct {
	foundEvent1 bool
	foundEvent2 bool
}

func (l *eventBusQueryListener) HandleEvent1(e eventBusTestEvent1) {
	l.foundEvent1 = true
}
func (l *eventBusQueryListener) HandleEvent2(e eventBusTestEvent2) {
	l.foundEvent2 = true
}
