package e2e

import (
	"errors"
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
	createFoo            = createFooCommand{fooId}
	nameFoo              = nameFooCommand{fooId, "a name"}
	notConfigured = notConfiguredCommand{}
	createBar            = createBarCommand{barId}
	configureBar         = configureBarCommand{barId, "a configuration"}
)

func TestCommandGateway_foo(t *testing.T) {
	eventBus := cqrs.NewEventBus()
	eventStore := cqrs.NewMemEventStore(eventBus)
	commandGateway := cqrs.NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&fooAggregate{})
	eventBus.RegisterQueryEventHandlers(&fooBarEventListener{})

	dispatchCleanly(commandGateway, createFoo)
	err := commandGateway.Dispatch(nameFoo)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(eventStore.Load(createFoo.Id)))
}

func TestCommandGateway_bar(t *testing.T) {
	eventBus := cqrs.NewEventBus()
	eventStore := cqrs.NewMemEventStore(eventBus)
	commandGateway := cqrs.NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&barAggregate{})
	eventBus.RegisterQueryEventHandlers(&fooBarEventListener{})

	dispatchCleanly(commandGateway, createBar)
	dispatchCleanly(commandGateway, configureBar)

	assert.Equal(t, 2, len(eventStore.Load(createBar.Id)))
}

func TestCombinedCommandGateways(t *testing.T) {
	eventBus := cqrs.NewEventBus()
	eventStore := cqrs.NewMemEventStore(eventBus)
	commandGateway := cqrs.NewCommandGateway(eventStore)
	commandGateway.RegisterAggregate(&fooAggregate{})
	commandGateway.RegisterAggregate(&barAggregate{})
	eventBus.RegisterQueryEventHandlers(&fooBarEventListener{})

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

type notConfiguredCommand struct {
	Id string
}

func (e notConfiguredCommand) TargetAggregateId() string { return e.Id }

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

var queryMap map[string]fooBarQuery

func init() {
	queryMap = make(map[string]fooBarQuery)
}

type fooBarQuery struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Configuration string `json:"configuration"`
}

type fooBarEventListener struct{}

func (*fooBarEventListener) OnBarCreated(e barCreatedEvent) {
	q := fooBarQuery{}
	q.Id = e.Id
	q.Type = "Bar"
	queryMap[e.Id] = q
}
func (*fooBarEventListener) OnBarConfigured(e barConfiguredEvent) {
	q := queryMap[e.Id]
	q.Configuration = e.Configuration
	queryMap[e.Id] = q
}
func (*fooBarEventListener) OnFooCreated(e fooCreatedEvent) {
	q := queryMap[e.Id]
	q.Id = e.Id
	q.Type = "Foo"
	queryMap[e.Id] = q
}
func (*fooBarEventListener) OnFooNamed(e fooNamedEvent) {
	q := queryMap[e.Id]
	q.Name = e.Name
	queryMap[e.Id] = q
}
