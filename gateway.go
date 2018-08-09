package cqrs

import (
	"reflect"
	"fmt"
)

type EventWrapper struct {
	AggregateId   string
	AggregateType reflect.Type
	Payload       interface{}
}

type CommandGateway struct {
	aggregateType   reflect.Type
	commandHandlers map[reflect.Type]*CommandHandler
}

func NewCommandGateway(aggregate interface{}) *CommandGateway {
	aggregateType := reflect.TypeOf(aggregate)
	return &CommandGateway{aggregateType, make(map[reflect.Type]*CommandHandler)}
}

func (gateway *CommandGateway) Register(aggregate interface{}) {
	aggregateType := reflect.TypeOf(aggregate)

	for i := 0; i < aggregateType.NumMethod(); i++ {
		f := aggregateType.Method(i)
		// TODO verify function signature before assuming that it is a command handler
		if f.Type.NumIn() != 2 {
			s := fmt.Sprintf("Currently only one argument commands handlers are supported, found %d in handler %s\n", f.Type.NumIn(), f.Name)
			panic(s)
		}

		in := f.Type.In(1)
		gateway.commandHandlers[in] = &CommandHandler{
			aggregateType: aggregateType,
			aggregate:     reflect.ValueOf(aggregate),
			commandType:   in,
			funcName:      f.Name,
			f:             f.Func,
		}

	}
	fmt.Println("Configured command handlers:")
	for k, v := range gateway.commandHandlers {
		fmt.Printf("\t%v - %s\n", k, v.funcName)
	}
}
func (gateway *CommandGateway) Dispatch(command interface{}) error {
	commandType := reflect.TypeOf(command)
	commandHandler := gateway.commandHandlers[commandType]
	if commandHandler == nil {
		s := fmt.Sprintf("Command handler for %v not configured", commandType)
		panic(s)
	}
	in := []reflect.Value{commandHandler.aggregate, reflect.ValueOf(command)}
	response := commandHandler.f.Call(in)
	//val := response[0]
	err := response[1].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}


func WrapEvent(aggregateId string, aggregateType reflect.Type, event interface{}) EventWrapper {
	return EventWrapper{
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Payload:       event,
	}
}

type CommandHandler struct {
	aggregateType reflect.Type
	aggregate     reflect.Value
	commandType   reflect.Type
	funcName      string
	f             reflect.Value
}
