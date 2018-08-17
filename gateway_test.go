package cqrs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	fooId = "a_foo_id"
	barId = "a_bar_id"
)

var (
	createFoo            = CreateFooCommand{fooId}
	nameFoo              = NameFooCommand{fooId, "a name"}
	notConfiguredCommand = NotConfiguredCommand{}
	createBar            = CreateBarCommand{barId}
	configureBar         = ConfigureBarCommand{barId, "a configuration"}
)

func TestCommandGateway_foo(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&FooAggregate{})

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.aggregateEventListeners))
}

func TestCommandGateway_errorOnDispatch(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&FooAggregate{})

	err := commandGateway.Dispatch(nameFoo)
	assert.NotNil(t, err)

	assert.Equal(t, 0, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_unconfiguredCommand(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&FooAggregate{})

	err := commandGateway.Dispatch(notConfiguredCommand)
	assert.NotNil(t, err)

	assert.Equal(t, 0, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_bar(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&BarAggregate{})

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.aggregateEventListeners))
}

func TestCombinedCommandGateways(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&FooAggregate{})
	commandGateway.RegisterAggregate(&BarAggregate{})

	assert.Equal(t, 4, len(commandGateway.commandHandlers))
	assert.Equal(t, 4, len(commandGateway.aggregateEventListeners))

}
