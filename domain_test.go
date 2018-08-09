package cqrs

import (
	"errors"
	"reflect"
)

type FooAggregate struct {
	id string
	name   string
}

type BarAggregate struct {
	someId        string
	configuration string
}

func (a *FooAggregate) HandleCreateFoo(e CreateFooCommand) ([]EventWrapper, error) {
	event := FooCreatedEvent{e.Id}
	eventWrapper := EventWrapper{
		AggregateId:   e.Id,
		AggregateType: reflect.TypeOf(a),
		Payload:       event,
	}
	return []EventWrapper{eventWrapper}, nil
}

func (a *FooAggregate) HandleNameFoo(e NameFooCommand) ([]EventWrapper, error) {
	if a.id == "" {
		return nil, errors.New("aggregate has not been initialized")
	}
	eventWrapper := WrapEvent(a.id, reflect.TypeOf(a), FooNamedEvent{e.Name})
	return []EventWrapper{eventWrapper}, nil
}

func (a *FooAggregate) onFooCreated(e *FooCreatedEvent) error {
	a.id = e.Id
	return nil
}

func (a *FooAggregate) onFooNamed(e *FooNamedEvent) error {
	a.name = e.Name
	return nil
}

type CreateFooCommand struct {
	Id string
}
type NameFooCommand struct {
	Name string
}
type NotConfiguredCommand struct {
}

type FooCreatedEvent struct {
	Id string
}
type FooNamedEvent struct {
	Name string
}

type CreateBarCommand struct {
	Id string
}
type ConfigureBarCommand struct {
	Name string
}

type BarCreatedEvent struct {
	Id string
}
type BarConfiguredEvent struct {
	Configuration string
}
