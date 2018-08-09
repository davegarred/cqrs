package cqrs

import (
	"reflect"
	"fmt"
)

var (
	domainObjectInterface = reflect.TypeOf((*DomainObject)(nil)).Elem()
	errorInterface  = reflect.TypeOf((*error)(nil)).Elem()
)

type EventWrapper struct {
	AggregateId   string
	AggregateType reflect.Type
	Payload       interface{}
}

type DomainObject interface {
	AggregateId() string
}

type CommandGateway struct {
	eventStore EventStore
	//aggregateType   reflect.Type
	commandHandlers map[reflect.Type]*CommandHandler
	eventListeners  map[reflect.Type]*CommandHandler
}

func NewCommandGateway(eventStore EventStore) *CommandGateway {
	return &CommandGateway{eventStore, make(map[reflect.Type]*CommandHandler), make(map[reflect.Type]*CommandHandler)}
}

func (gateway *CommandGateway) Register(aggregate interface{}) {
	aggregateType := reflect.TypeOf(aggregate)

	for i := 0; i < aggregateType.NumMethod(); i++ {
		f := aggregateType.Method(i)

		if hasCommandHandlerSignature(f) {
			in := f.Type.In(1)

				gateway.commandHandlers[in] = &CommandHandler{
					aggregateType: aggregateType,
					aggregate:     reflect.ValueOf(aggregate),
					funcName:      f.Name,
					f:             f.Func,
				}
		} else if hasEventListenerSignature(f) {
			in := f.Type.In(1)
			gateway.eventListeners[in] = &CommandHandler{
				aggregateType: aggregateType,
				aggregate:     reflect.ValueOf(aggregate),
				funcName:      f.Name,
				f:             f.Func,
			}
		}
	}
	fmt.Printf("Aggregate %v registered\n", aggregateType)
	fmt.Println("Configured command handlers:")
	for k, v := range gateway.commandHandlers {
		fmt.Printf("\t%v - %s\n", k, v.funcName)
	}
	fmt.Println("Configured event listeners:")
	for k, v := range gateway.eventListeners {
		fmt.Printf("\t%v - %s\n", k, v.funcName)
	}
}
func hasCommandHandlerSignature(f reflect.Method) bool {
	if f.Type.NumIn() != 2 || f.Type.NumOut() != 2 {
		return false
	}

	if !f.Type.In(1).Implements(domainObjectInterface) {
		return false
	}

	var wrap []EventWrapper

	return f.Type.Out(0) == reflect.TypeOf(wrap) && f.Type.Out(1).Implements(errorInterface)
}
func hasEventListenerSignature(f reflect.Method) bool {
	if f.Type.NumIn() != 2 || f.Type.NumOut() > 0 {
		return false
	}
	return f.Type.In(1).Implements(domainObjectInterface)
}

func (gateway *CommandGateway) Dispatch(command DomainObject) error {
	commandType := reflect.TypeOf(command)
	commandHandler := gateway.commandHandlers[commandType]
	if commandHandler == nil {
		s := fmt.Sprintf("DomainObject handler for %v not configured", commandType)
		panic(s)
	}

	aggregateId := command.AggregateId()
	aggregate := gateway.loadAggregate(commandHandler.aggregateType, aggregateId)
	in := []reflect.Value{aggregate, reflect.ValueOf(command)}

	response := commandHandler.f.Call(in)
	err := response[1].Interface()
	if err != nil {
		return err.(error)
	}
	eventWrappers := response[0].Interface().([]EventWrapper)
	gateway.eventStore.Persist(aggregateId, eventWrappers)
	return nil
}
func (gateway *CommandGateway) loadAggregate(aggregateType reflect.Type, aggregateId string) reflect.Value {
	events := gateway.eventStore.Load(aggregateId)
	aggregate := reflect.New(aggregateType.Elem())
	for _,event := range events {
		fmt.Println("applying event: ", event)
	}
	return aggregate
}

func WrapEvent(aggregateId string, aggregate interface{}, event interface{}) EventWrapper {
	return EventWrapper{
		AggregateId:   aggregateId,
		AggregateType: reflect.TypeOf(aggregate),
		Payload:       event,
	}
}

type CommandHandler struct {
	aggregateType reflect.Type
	aggregate     reflect.Value
	funcName      string
	f             reflect.Value
}
