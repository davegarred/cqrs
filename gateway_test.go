package cqrs

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&FooAggregate{})
	commandGateway.RegisterQueryEventHandlers(&FooBarEventListener{})

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.eventListeners))
	assert.Equal(t, 4, len(commandGateway.queryEventListeners))

	dispatchCleanly(commandGateway, createFoo)
	err := commandGateway.Dispatch(nameFoo)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_errorOnDispatch(t *testing.T) {
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&FooAggregate{})

	err := commandGateway.Dispatch(nameFoo)
	assert.NotNil(t, err)

	assert.Equal(t, 0, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_unconfiguredCommand(t *testing.T) {
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&FooAggregate{})

	err := commandGateway.Dispatch(notConfiguredCommand)
	assert.NotNil(t, err)

	assert.Equal(t, 0, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_bar(t *testing.T) {
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&BarAggregate{})
	commandGateway.RegisterQueryEventHandlers(&FooBarEventListener{})

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.eventListeners))
	assert.Equal(t, 4, len(commandGateway.queryEventListeners))

	dispatchCleanly(commandGateway, createBar)
	dispatchCleanly(commandGateway, configureBar)

	assert.Equal(t, 2, len(eventStore.Load(createBar.Id)))
}


func TestCombinedCommandGateways(t *testing.T) {
	eventStore := NewMemEventStore()
	commandGateway := NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&FooAggregate{})
	commandGateway.RegisterAggregate(&BarAggregate{})
	commandGateway.RegisterQueryEventHandlers(&FooBarEventListener{})

	assert.Equal(t, 4, len(commandGateway.commandHandlers))
	assert.Equal(t, 4, len(commandGateway.eventListeners))
	assert.Equal(t, 4, len(commandGateway.queryEventListeners))

	dispatchCleanly(commandGateway, createFoo)
	dispatchCleanly(commandGateway, nameFoo)

	dispatchCleanly(commandGateway, createBar)
	dispatchCleanly(commandGateway, configureBar)


	fmt.Println("Published events:")
	for _,event := range eventStore.Load(createBar.Id) {
		fmt.Printf("\t- %+v\n", event)
	}
	for _,event := range eventStore.Load(createFoo.Id) {
		fmt.Printf("\t- %+v\n", event)
	}
}

func dispatchCleanly(commandGateway *CommandGateway, c Command) error {
	err := commandGateway.Dispatch(c)
	if err != nil {
		panic(err)
	}
	return nil
}