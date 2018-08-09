package cqrs

import (
	"testing"
	"fmt"
)

var (
	createFoo   = CreateFooCommand{"an id"}
	nameFoo     = NameFooCommand{"a name"}
	genericSome = NotConfiguredCommand{}
)

func TestCommandGateway(t *testing.T) {
	fooAggregate := &FooAggregate{}
	commandGateway := NewCommandGateway(fooAggregate)
	commandGateway.Register(fooAggregate)

	dispatch(commandGateway, createFoo)
	//dispatch(commandGateway, nameFoo)

	fmt.Printf("Final aggregate state: %+v\n", *fooAggregate)
}

func dispatch(commandGateway *CommandGateway, c interface{}) error {
	err := commandGateway.Dispatch(c)
	if err != nil {
		panic(err)
	}
	return nil
}
