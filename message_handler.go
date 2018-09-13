package cqrs

import (
	"github.com/davegarred/cqrs/ext"
	"reflect"
)

var (
	commandInterface    = reflect.TypeOf((*ext.Command)(nil)).Elem()
	eventInterface      = reflect.TypeOf((*ext.Event)(nil)).Elem()
	eventSliceInterface = reflect.TypeOf([]ext.Event{})
	errorInterface      = reflect.TypeOf((*error)(nil)).Elem()
)

type aggregateMessageHandler struct {
	AggregateType reflect.Type
	FuncName      string
	F             reflect.Value
}

func NewMessageHandler(aggregateType reflect.Type, f reflect.Method) *aggregateMessageHandler {
	return &aggregateMessageHandler{
		AggregateType: aggregateType,
		FuncName:      f.Name,
		F:             f.Func,
	}
}

type queryEventListener struct {
	Query    interface{}
	FuncName string
	F        reflect.Value
}

func NewEventListener(query interface{}, f reflect.Method) *queryEventListener {
	return &queryEventListener{
		Query:    query,
		FuncName: f.Name,
		F:        f.Func,
	}
}

func (handler *aggregateMessageHandler) applyCommand(aggregate reflect.Value, command ext.Command) ([]ext.Event, error) {
	in := []reflect.Value{aggregate, reflect.ValueOf(command)}
	response := handler.F.Call(in)
	err := response[1].Interface()
	if err != nil {
		return nil, err.(error)
	}
	events := response[0].Interface().([]ext.Event)
	return events, nil
}

func (handler *aggregateMessageHandler) applyEvent(aggregate reflect.Value, event ext.Event) {
	in := []reflect.Value{aggregate, reflect.ValueOf(event)}
	handler.F.Call(in)
}

func (handler *queryEventListener) applyEvent(event ext.Event) {
	in := []reflect.Value{reflect.ValueOf(handler.Query), reflect.ValueOf(event)}
	handler.F.Call(in)
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
