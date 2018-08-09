package cqrs

import (
	"errors"
	"reflect"
)

type FooAggregate struct {
	fooId string
	name  string
}


func (a *FooAggregate) NonCQRSFunction_noParams() {
	panic("This should never be called")
}
func (a *FooAggregate) NonCQRSFunction_oneParam(_ string) {
	panic("This should never be called")
}
func (a *FooAggregate) NonCQRSFunction_twoParams(_ string, _ int) {
	panic("This should never be called")
}
func (a *FooAggregate) NonCQRSFunction_oneParam_similarSig(_ string) ([]string,error) {
	panic("This should never be called")
	return nil, nil
}
func (a *FooAggregate) HandleCreateFoo(e CreateFooCommand) ([]EventWrapper, error) {
	eventWrapper := WrapEvent(e.Id, a, FooCreatedEvent{e.Id})
	return []EventWrapper{eventWrapper}, nil
}
func (a *FooAggregate) HandleNameFoo_requireIdSet(e NameFooCommand) ([]EventWrapper, error) {
	if a.fooId == "" {
		return nil, errors.New("aggregate has not been initialized")
	}
	eventWrapper := WrapEvent(a.fooId, a, FooNamedEvent{e.Id, e.Name})
	return []EventWrapper{eventWrapper}, nil
}

func (a *FooAggregate) OnFooCreated(e *FooCreatedEvent) {
	a.fooId = e.Id
}
func (a *FooAggregate) OnFooNamed(e *FooNamedEvent) {
	a.name = e.Name
}

type CreateFooCommand struct {
	Id string
}
func (e CreateFooCommand) AggregateId() string {return e.Id}
type NameFooCommand struct {
	Id string
	Name string
}
func (e NameFooCommand) AggregateId() string {return e.Id}
type NotConfiguredCommand struct {
	Id string
}
func (e NotConfiguredCommand) AggregateId() string {return e.Id}

type FooCreatedEvent struct {
	Id string
}
func (e FooCreatedEvent) AggregateId() string {return e.Id}
type FooNamedEvent struct {
	Id string
	Name string
}
func (e FooNamedEvent) AggregateId() string {return e.Id}



type BarAggregate struct {
	barId         string
	configuration string
}

func (a *BarAggregate) HandleCreateBar(e CreateBarCommand) ([]EventWrapper, error) {
	event := FooCreatedEvent{e.Id}
	eventWrapper := WrapEvent(e.Id,reflect.TypeOf(a), event)
	return []EventWrapper{eventWrapper}, nil
}
func (a *BarAggregate) HandleNameBar(e ConfigureBarCommand) ([]EventWrapper, error) {
	eventWrapper := WrapEvent(a.barId, a, BarConfiguredEvent{e.Id, e.Configuration})
	return []EventWrapper{eventWrapper}, nil
}

func (a *BarAggregate) OnBarCreated(e BarCreatedEvent) {
	a.barId = e.Id
}
func (a *BarAggregate) OnBarConfigured(e BarConfiguredEvent) {
	a.configuration = e.Configuration
}

type CreateBarCommand struct {
	Id string
}
func (e CreateBarCommand) AggregateId() string {return e.Id}
type ConfigureBarCommand struct {
	Id string
	Configuration string
}
func (e ConfigureBarCommand) AggregateId() string {return e.Id}

type BarCreatedEvent struct {
	Id string
}
func (e BarCreatedEvent) AggregateId() string {return e.Id}
type BarConfiguredEvent struct {
	Id string
	Configuration string
}
func (e BarConfiguredEvent) AggregateId() string {return e.Id}
