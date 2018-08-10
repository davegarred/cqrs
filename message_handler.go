package cqrs

import (
	"reflect"
)

type MessageHandler struct {
	aggregateType reflect.Type
	aggregate     reflect.Value
	funcName      string
	f             reflect.Value
}

func NewMessageHandler(aggregateType reflect.Type, f reflect.Method) *MessageHandler {
	return &MessageHandler{
		aggregateType: aggregateType,
		//aggregate:     aggregateValue,
		funcName:      f.Name,
		f:             f.Func,
	}
}

func (handler *MessageHandler) Call(aggregate reflect.Value, command Command) ([]Event, error) {
	in := []reflect.Value{aggregate, reflect.ValueOf(command)}
	response := handler.f.Call(in)
	err := response[1].Interface()
	if err != nil {
		return nil, err.(error)
	}
	events := response[0].Interface().([]Event)
	return events, nil
}

func (handler *MessageHandler) ApplyEvent(aggregate reflect.Value, event Event) reflect.Value {
	in := []reflect.Value{aggregate, reflect.ValueOf(event)}
	handler.f.Call(in)
	return aggregate
}

func hasCommandHandlerSignature(f reflect.Method) bool {
	if f.Type.NumIn() != 2 || f.Type.NumOut() != 2 {
		return false
	}
	takesCommand := f.Type.In(1).Implements(commandInterface)
	returnsEventsFirst := f.Type.Out(0) == eventSliceInterface
	returnsErrorSecond := f.Type.Out(1).Implements(errorInterface)
	return takesCommand && returnsEventsFirst && returnsErrorSecond
}
func hasEventListenerSignature(f reflect.Method) bool {
	if f.Type.NumIn() != 2 || f.Type.NumOut() > 0 {
		return false
	}
	takesEvent := f.Type.In(1).Implements(eventInterface)
	return takesEvent
}
