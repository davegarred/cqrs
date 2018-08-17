package test

import (
	"fmt"
	"github.com/davegarred/cqrs"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	fooId = "a_foo_id"
	barId = "a_bar_id"
)

var (
	createFoo            = cqrs.CreateFooCommand{fooId}
	nameFoo              = cqrs.NameFooCommand{fooId, "a name"}
	notConfiguredCommand = cqrs.NotConfiguredCommand{}
	createBar            = cqrs.CreateBarCommand{barId}
	configureBar         = cqrs.ConfigureBarCommand{barId, "a configuration"}
)

func TestCommandGateway_foo(t *testing.T) {
	eventBus := cqrs.NewEventBus()
	eventStore := cqrs.NewMemEventStore()
	commandGateway := cqrs.NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&cqrs.FooAggregate{})
	eventBus.RegisterQueryEventHandlers(&cqrs.FooBarEventListener{})

	dispatchCleanly(commandGateway, createFoo)
	err := commandGateway.Dispatch(nameFoo)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_bar(t *testing.T) {
	eventBus := cqrs.NewEventBus()
	eventStore := cqrs.NewMemEventStore()
	commandGateway := cqrs.NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&cqrs.BarAggregate{})
	eventBus.RegisterQueryEventHandlers(&cqrs.FooBarEventListener{})

	dispatchCleanly(commandGateway, createBar)
	dispatchCleanly(commandGateway, configureBar)

	assert.Equal(t, 2, len(eventStore.Load(createBar.Id)))
}

func TestCombinedCommandGateways(t *testing.T) {
	eventBus := cqrs.NewEventBus()
	eventStore := cqrs.NewMemEventStore()
	commandGateway := cqrs.NewCommandGateway(eventStore, eventBus)
	commandGateway.RegisterAggregate(&cqrs.FooAggregate{})
	commandGateway.RegisterAggregate(&cqrs.BarAggregate{})
	eventBus.RegisterQueryEventHandlers(&cqrs.FooBarEventListener{})

	dispatchCleanly(commandGateway, createFoo)
	dispatchCleanly(commandGateway, nameFoo)

	dispatchCleanly(commandGateway, createBar)
	dispatchCleanly(commandGateway, configureBar)

	fmt.Println("Published events:")
	for _, event := range eventStore.Load(createBar.Id) {
		fmt.Printf("\t- %+v\n", event)
	}
	for _, event := range eventStore.Load(createFoo.Id) {
		fmt.Printf("\t- %+v\n", event)
	}
}

func dispatchCleanly(commandGateway *cqrs.CommandGateway, c cqrs.Command) error {
	err := commandGateway.Dispatch(c)
	if err != nil {
		panic(err)
	}
	return nil
}
