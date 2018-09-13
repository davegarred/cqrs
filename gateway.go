package cqrs

import (
	"errors"
	"fmt"
	"github.com/davegarred/cqrs/ext"
	"reflect"
)

type CommandGateway struct {
	eventStore              ext.EventStore
	commandHandlers         map[reflect.Type]*aggregateMessageHandler
	aggregateEventListeners map[reflect.Type]*aggregateMessageHandler
}

func NewCommandGateway(eventStore ext.EventStore) *CommandGateway {
	return &CommandGateway{eventStore, make(map[reflect.Type]*aggregateMessageHandler), make(map[reflect.Type]*aggregateMessageHandler)}
}

func (gateway *CommandGateway) RegisterAggregate(aggregate interface{}) {
	aggregateType := reflect.TypeOf(aggregate)

	for i := 0; i < aggregateType.NumMethod(); i++ {
		f := aggregateType.Method(i)

		if hasCommandHandlerSignature(f) {
			gateway.commandHandlers[f.Type.In(1)] = NewMessageHandler(aggregateType, f)
		} else if hasEventListenerSignature(f) {
			gateway.aggregateEventListeners[f.Type.In(1)] = NewMessageHandler(aggregateType, f)
		}
	}
}

func (gateway *CommandGateway) Dispatch(command ext.Command) error {
	commandType := reflect.TypeOf(command)
	commandHandler := gateway.commandHandlers[commandType]
	if commandHandler == nil {
		s := fmt.Sprintf("Command handler for %v not configured", commandType)
		return errors.New(s)
	}

	aggregateId := command.TargetAggregateId()
	aggregate := gateway.loadAggregate(commandHandler.AggregateType, aggregateId)

	events, err := commandHandler.applyCommand(aggregate, command)
	if err != nil {
		return err
	}
	gateway.eventStore.Persist(aggregateId, events)
	return nil
}

func (gateway *CommandGateway) loadAggregate(aggregateType reflect.Type, aggregateId string) reflect.Value {
	events := gateway.eventStore.Load(aggregateId)
	aggregate := reflect.New(aggregateType.Elem())
	for _, event := range events {
		listener := gateway.aggregateEventListeners[reflect.TypeOf(event)]
		if listener != nil {
			if listener.AggregateType != aggregateType {
				error := fmt.Sprintf("Incorrectly configured event listener, event type %T was produced via %v but has an event listener attached to %v\n", event, aggregateType, listener.AggregateType)
				panic(error)
			}
			listener.applyEvent(aggregate, event)
		}
	}
	return aggregate
}
