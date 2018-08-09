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
	createFoo      = CreateFooCommand{fooId}
	nameFoo        = NameFooCommand{fooId, "a name"}
	genericCommand = NotConfiguredCommand{}
	createBar      = CreateBarCommand{barId}
	configureBar   = ConfigureBarCommand{barId, "a configuration"}
)

func TestCommandGateway_foo(t *testing.T) {
	fooAggregate := &FooAggregate{}
	commandGateway := NewCommandGateway(NewMemEventStore())
	commandGateway.Register(fooAggregate)

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.eventListeners))

	dispatchCleanly(commandGateway, createFoo)
	err := commandGateway.Dispatch(nameFoo)
	assert.NotNil(t, err)

	fmt.Printf("Final aggregate state: %+v\n", *fooAggregate)
}

func TestCommandGateway_bar(t *testing.T) {
	barAggregate := &BarAggregate{}
	commandGateway := NewCommandGateway(NewMemEventStore())
	commandGateway.Register(barAggregate)

	assert.Equal(t, 2, len(commandGateway.commandHandlers))
	assert.Equal(t, 2, len(commandGateway.eventListeners))

	dispatchCleanly(commandGateway, createBar)
	dispatchCleanly(commandGateway, configureBar)

	fmt.Printf("Final aggregate state: %+v\n", *barAggregate)
}

func dispatchCleanly(commandGateway *CommandGateway, c DomainObject) error {
	err := commandGateway.Dispatch(c)
	if err != nil {
		panic(err)
	}
	return nil
}