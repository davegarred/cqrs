package components

import (
	"errors"
	"github.com/davegarred/cqrs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventBus_RegisterQueryEventHandlers(t *testing.T) {
	eventBus := NewEventBus()
	eventBus.RegisterQueryEventHandlers(&fooBarEventListener{})

	assert.Equal(t, 4, len(eventBus.queryEventListeners))
}

type fooBarEventListener struct {}

func (*fooBarEventListener) OnBarCreated(e barCreatedEvent)       {}
func (*fooBarEventListener) OnBarConfigured(e barConfiguredEvent) {}
func (*fooBarEventListener) OnFooCreated(e fooCreatedEvent)       {}
func (*fooBarEventListener) OnFooNamed(e fooNamedEvent)           {}

type fooAggregate struct {
	fooId string
	name  string
}

func (a *fooAggregate) NonCQRSFunction_noParams() {
	panic("This should never be called")
}
func (a *fooAggregate) NonCQRSFunction_oneParam(_ string) {
	panic("This should never be called")
}
func (a *fooAggregate) NonCQRSFunction_twoParams(_ string, _ int) {
	panic("This should never be called")
}
func (a *fooAggregate) NonCQRSFunction_oneParam_similarSig(_ string) ([]string, error) {
	panic("This should never be called")
	return nil, nil
}
func (a *fooAggregate) HandleCreateFoo(e createFooCommand) ([]cqrs.Event, error) {
	return []cqrs.Event{fooCreatedEvent{e.Id}}, nil
}
func (a *fooAggregate) HandleNameFoo_requireIdSet(e nameFooCommand) ([]cqrs.Event, error) {
	if a.fooId == "" {
		return nil, errors.New("aggregate has not been initialized")
	}
	return []cqrs.Event{fooNamedEvent{e.Id, e.Name}}, nil
}

func (a *fooAggregate) OnFooCreated(e fooCreatedEvent) {
	a.fooId = e.Id
}
func (a *fooAggregate) OnFooNamed(e fooNamedEvent) {
	a.name = e.Name
}

type createFooCommand struct {
	Id string
}

func (e createFooCommand) TargetAggregateId() string { return e.Id }

type nameFooCommand struct {
	Id   string
	Name string
}

func (e nameFooCommand) TargetAggregateId() string { return e.Id }

type fooCreatedEvent struct {
	Id string
}

func (e fooCreatedEvent) AggregateId() string { return e.Id }

type fooNamedEvent struct {
	Id   string
	Name string
}

func (e fooNamedEvent) AggregateId() string { return e.Id }

type barAggregate struct {
	barId         string
	configuration string
}

func (a *barAggregate) HandleCreateBar(e createBarCommand) ([]cqrs.Event, error) {
	return []cqrs.Event{barCreatedEvent{e.Id}}, nil
}
func (a *barAggregate) HandleNameBar(e configureBarCommand) ([]cqrs.Event, error) {
	return []cqrs.Event{barConfiguredEvent{e.Id, e.Configuration}}, nil
}

func (a *barAggregate) OnBarCreated(e barCreatedEvent) {
	a.barId = e.Id
}
func (a *barAggregate) OnBarConfigured(e barConfiguredEvent) {
	a.configuration = e.Configuration
}

type createBarCommand struct {
	Id string
}

func (e createBarCommand) TargetAggregateId() string { return e.Id }

type configureBarCommand struct {
	Id            string
	Configuration string
}

func (e configureBarCommand) TargetAggregateId() string { return e.Id }

type barCreatedEvent struct {
	Id string
}

func (e barCreatedEvent) AggregateId() string { return e.Id }

type barConfiguredEvent struct {
	Id            string
	Configuration string
}

func (e barConfiguredEvent) AggregateId() string { return e.Id }
