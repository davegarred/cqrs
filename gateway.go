package cqrs

import (
	"errors"
	"fmt"
	"reflect"
)

type CommandGateway struct {
	eventStore              EventStore
	eventBus                *EventBus
	commandHandlers         map[reflect.Type]*MessageHandler
	aggregateEventListeners map[reflect.Type]*MessageHandler
}

func NewCommandGateway(eventStore EventStore, eventBus *EventBus) *CommandGateway {
	return &CommandGateway{eventStore, eventBus, make(map[reflect.Type]*MessageHandler), make(map[reflect.Type]*MessageHandler)}
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

func (gateway *CommandGateway) Dispatch(command Command) error {
	commandType := reflect.TypeOf(command)
	commandHandler := gateway.commandHandlers[commandType]
	if commandHandler == nil {
		s := fmt.Sprintf("Command handler for %v not configured", commandType)
		return errors.New(s)
	}

	aggregateId := command.TargetAggregateId()
	aggregate := gateway.loadAggregate(commandHandler.AggregateType, aggregateId)

	events, err := commandHandler.Call(aggregate, command)
	if err != nil {
		return err
	}
	gateway.eventStore.Persist(aggregateId, events)
	gateway.eventBus.PublishEvents(events)
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
			listener.ApplyEvent(aggregate, event)
		}
	}
	return aggregate
}
