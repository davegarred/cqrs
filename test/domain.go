package test

import (
	"errors"
	"github.com/davegarred/cqrs"
)

type FooAggregate struct {
	fooId string
	name  string
}

func (a *FooAggregate) HandleCreateFoo(e CreateFooCommand) ([]cqrs.Event, error) {
	return []cqrs.Event{FooCreatedEvent{e.Id}}, nil
}
func (a *FooAggregate) HandleNameFoo_requireIdSet(e NameFooCommand) ([]cqrs.Event, error) {
	if a.fooId == "" {
		return nil, errors.New("aggregate has not been initialized")
	}
	return []cqrs.Event{FooNamedEvent{e.Id, e.Name}}, nil
}

func (a *FooAggregate) OnFooCreated(e FooCreatedEvent) {
	a.fooId = e.Id
}
func (a *FooAggregate) OnFooNamed(e FooNamedEvent) {
	a.name = e.Name
}

type CreateFooCommand struct {
	Id string
}

func (e CreateFooCommand) TargetAggregateId() string { return e.Id }

type NameFooCommand struct {
	Id   string
	Name string
}

func (e NameFooCommand) TargetAggregateId() string { return e.Id }

type NotConfiguredCommand struct {
	Id string
}

func (e NotConfiguredCommand) TargetAggregateId() string { return e.Id }

type FooCreatedEvent struct {
	Id string
}

func (e FooCreatedEvent) AggregateId() string { return e.Id }

type FooNamedEvent struct {
	Id   string
	Name string
}

func (e FooNamedEvent) AggregateId() string { return e.Id }

type BarAggregate struct {
	barId         string
	configuration string
}

func (a *BarAggregate) HandleCreateBar(e CreateBarCommand) ([]cqrs.Event, error) {
	return []cqrs.Event{BarCreatedEvent{e.Id}}, nil
}
func (a *BarAggregate) HandleNameBar(e ConfigureBarCommand) ([]cqrs.Event, error) {
	return []cqrs.Event{BarConfiguredEvent{e.Id, e.Configuration}}, nil
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

func (e CreateBarCommand) TargetAggregateId() string { return e.Id }

type ConfigureBarCommand struct {
	Id            string
	Configuration string
}

func (e ConfigureBarCommand) TargetAggregateId() string { return e.Id }

type BarCreatedEvent struct {
	Id string
}

func (e BarCreatedEvent) AggregateId() string { return e.Id }

type BarConfiguredEvent struct {
	Id            string
	Configuration string
}

func (e BarConfiguredEvent) AggregateId() string { return e.Id }

var queryMap map[string]FooBarQuery

func init() {
	queryMap = make(map[string]FooBarQuery)
}

type FooBarQuery struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Configuration string `json:"configuration"`
}

type FooBarEventListener cqrs.QueryEventListener

func (*FooBarEventListener) OnBarCreated(e BarCreatedEvent) {
	q := FooBarQuery{}
	q.Id = e.Id
	q.Type = "Bar"
	queryMap[e.Id] = q
}
func (*FooBarEventListener) OnBarConfigured(e BarConfiguredEvent) {
	q := queryMap[e.Id]
	q.Configuration = e.Configuration
	queryMap[e.Id] = q
}
func (*FooBarEventListener) OnFooCreated(e FooCreatedEvent) {
	q := queryMap[e.Id]
	q.Id = e.Id
	q.Type = "Foo"
	queryMap[e.Id] = q
}
func (*FooBarEventListener) OnFooNamed(e FooNamedEvent) {
	q := queryMap[e.Id]
	q.Name = e.Name
	queryMap[e.Id] = q
}
