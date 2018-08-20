package cqrs

import (
	"fmt"
	"reflect"
	"errors"
)

var (
	commandInterface    = reflect.TypeOf((*Command)(nil)).Elem()
	eventInterface      = reflect.TypeOf((*Event)(nil)).Elem()
	eventSliceInterface = reflect.TypeOf([]Event{})
	errorInterface      = reflect.TypeOf((*error)(nil)).Elem()
)

type Command interface {
	TargetAggregateId() string
}

type Event interface {
	AggregateId() string
}

type CommandGateway struct {
	eventStore      EventStore
	commandHandlers map[reflect.Type]*MessageHandler
	eventListeners  map[reflect.Type]*MessageHandler
	queryEventListeners  map[reflect.Type][]*MessageHandler
}

func NewCommandGateway(eventStore EventStore) *CommandGateway {
	return &CommandGateway{eventStore, make(map[reflect.Type]*MessageHandler), make(map[reflect.Type]*MessageHandler), make(map[reflect.Type][]*MessageHandler)}
}

func (gateway *CommandGateway) RegisterAggregate(aggregate interface{}) {
	aggregateType := reflect.TypeOf(aggregate)

	for i := 0; i < aggregateType.NumMethod(); i++ {
		f := aggregateType.Method(i)

		if hasCommandHandlerSignature(f) {
			gateway.commandHandlers[f.Type.In(1)]  = NewMessageHandler(aggregateType, f)
		} else if hasEventListenerSignature(f) {
			gateway.eventListeners[f.Type.In(1)] = NewMessageHandler(aggregateType, f)
		}
	}
}

func (gateway *CommandGateway) RegisterQueryEventHandlers(listener interface{}) {
	aggregateType := reflect.TypeOf(listener)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		f := aggregateType.Method(i)

		if hasEventListenerSignature(f) {
			eventType := f.Type.In(1)
			queryEventListeners := gateway.queryEventListeners[eventType]
			if queryEventListeners == nil {
				queryEventListeners = make([]*MessageHandler,0)
			}
			queryEventListeners = append(queryEventListeners, NewMessageHandler(aggregateType, f))
			gateway.queryEventListeners[eventType] = queryEventListeners
		}
	}
}

func (gateway *CommandGateway) logAggregateRegistrationDetails() {
	fmt.Println("Configured command handlers:")
	for k, v := range gateway.commandHandlers {
		fmt.Printf("\t%v - %s (%v)\n", k, v.funcName, v.aggregateType)
	}
	fmt.Println("Configured event listeners:")
	for k, v := range gateway.eventListeners {
		fmt.Printf("\t%v - %s (%v)\n", k, v.funcName, v.aggregateType)
	}
}

func (gateway *CommandGateway) Dispatch(command Command) error {
	commandType := reflect.TypeOf(command)
	commandHandler := gateway.commandHandlers[commandType]
	if commandHandler == nil {
		s := fmt.Sprintf("Command handler for %v not configured", commandType)
		return errors.New(s)
	}

	aggregateId := command.TargetAggregateId()
	aggregate := gateway.loadAggregate(commandHandler.aggregateType, aggregateId)

	events, err := commandHandler.Call(aggregate, command)
	if err != nil {
		return err
	}
	gateway.eventStore.Persist(aggregateId, events)
	publishEvents(events, gateway)
	return nil
}

func publishEvents(events []Event, gateway *CommandGateway) {
	for _, event := range events {
		for _,listener := range gateway.queryEventListeners[reflect.TypeOf(event)] {
			agg := reflect.New(listener.aggregateType).Elem()
			listener.ApplyEvent(agg, event)
		}
	}
}
func (gateway *CommandGateway) loadAggregate(aggregateType reflect.Type, aggregateId string) reflect.Value {
	events := gateway.eventStore.Load(aggregateId)
	aggregate := reflect.New(aggregateType.Elem())
	for _, event := range events {
		listener := gateway.eventListeners[reflect.TypeOf(event)]
		if listener != nil {
			if listener.aggregateType != aggregateType {
				error := fmt.Sprintf("Incorrectly configured event listener, event type %T was produced via %v but has an event listener attached to %v\n", event, aggregateType, listener.aggregateType)
				panic(error)
			}
			listener.ApplyEvent(aggregate, event)
		}
	}
	return aggregate
}
