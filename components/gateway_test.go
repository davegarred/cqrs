package components

import (
	"github.com/davegarred/cqrs/persist"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fooId = "a_foo_id"
	barId = "a_bar_id"
)

var (
	createFoo    = createFooCommand{fooId}
	nameFoo      = nameFooCommand{fooId, "a name"}
	createBar    = createBarCommand{barId}
	configureBar = configureBarCommand{barId, "a configuration"}
)

func TestCommandGateway_foo(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := persist.NewMemEventStore(eventBus)
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&fooAggregate{})

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.aggregateEventListeners))

	err := commandGateway.Dispatch(createFoo)
	assert.Nil(t, err)
}

func TestCommandGateway_errorOnDispatch(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := persist.NewMemEventStore(eventBus)
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&fooAggregate{})

	err := commandGateway.Dispatch(nameFoo)
	assert.NotNil(t, err)

	assert.Equal(t, 0, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_unconfiguredCommand(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := persist.NewMemEventStore(eventBus)
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&fooAggregate{})

	err := commandGateway.Dispatch(notConfiguredCommand{})
	assert.NotNil(t, err)

	assert.Equal(t, 0, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_bar(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := persist.NewMemEventStore(eventBus)
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&barAggregate{})

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.aggregateEventListeners))

	err := commandGateway.Dispatch(createBar)
	assert.Nil(t, err)
	err = commandGateway.Dispatch(configureBar)
	assert.Nil(t, err)
	err = commandGateway.Dispatch(createFoo)
	assert.NotNil(t, err)
}

func TestCombinedCommandGateways(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := persist.NewMemEventStore(eventBus)
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&fooAggregate{})
	commandGateway.RegisterAggregate(&barAggregate{})

	assert.Equal(t, 4, len(commandGateway.commandHandlers))
	assert.Equal(t, 4, len(commandGateway.aggregateEventListeners))
}

type notConfiguredCommand struct {
	Id string
}

func (e notConfiguredCommand) TargetAggregateId() string { return e.Id }
