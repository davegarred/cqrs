package cqrs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hasEventListenerSignature(t *testing.T) {
	eventListener, _ := reflect.TypeOf(&testMessageHandlerQueryEventListener{}).MethodByName("Handle")
	commandHandler, _ := reflect.TypeOf(&testMessageHandlerAggregate{}).MethodByName("Handle")
	aggregateEventListener, _ := reflect.TypeOf(&testMessageHandlerAggregate{}).MethodByName("HandleEvent")

	assert.True(t, hasEventListenerSignature(eventListener))
	assert.False(t, hasEventListenerSignature(commandHandler))
	assert.True(t, hasEventListenerSignature(aggregateEventListener))
}

func Test_hasCommandHandlerSignature(t *testing.T) {
	eventListener, _ := reflect.TypeOf(&testMessageHandlerQueryEventListener{}).MethodByName("Handle")
	commandHandler, _ := reflect.TypeOf(&testMessageHandlerAggregate{}).MethodByName("Handle")
	aggregateEventListener, _ := reflect.TypeOf(&testMessageHandlerAggregate{}).MethodByName("HandleEvent")

	assert.True(t, hasCommandHandlerSignature(commandHandler))
	assert.False(t, hasCommandHandlerSignature(eventListener))
	assert.False(t, hasCommandHandlerSignature(aggregateEventListener))
}

func Test_queryEventListener_applyEvent(t *testing.T) {
	listener := &testMessageHandlerQueryEventListener{}
	method, _ := reflect.TypeOf(listener).MethodByName("Handle")
	eventListener := NewEventListener(listener, method)

	eventListener.applyEvent(testMessageHandlerEvent{})

	assert.True(t, listener.success)
}

func Test_aggregateMessageHandler_applyCommand(t *testing.T) {
	aggregate := &testMessageHandlerAggregate{}
	aggregateType := reflect.TypeOf(aggregate)
	method, _ := aggregateType.MethodByName("Handle")
	messageHandler := NewMessageHandler(aggregateType, method)

	events, err := messageHandler.applyCommand(reflect.ValueOf(aggregate), testMessageHandlerCommand{})

	assert.Equal(t, []Event{testMessageHandlerEvent{}}, events)
	assert.Nil(t, err)
}

func Test_aggregateMessageHandler_applyEvent(t *testing.T) {
	aggregate := &testMessageHandlerAggregate{}
	aggregateType := reflect.TypeOf(aggregate)
	method, _ := aggregateType.MethodByName("HandleEvent")
	messageHandler := NewMessageHandler(aggregateType, method)

	messageHandler.applyEvent(reflect.ValueOf(aggregate), testMessageHandlerEvent{})

	assert.True(t, aggregate.success)
}

type testMessageHandlerEvent struct{}

func (e testMessageHandlerEvent) AggregateId() string { return "" }

type testMessageHandlerQueryEventListener struct {
	success bool
}

func (l *testMessageHandlerQueryEventListener) Handle(e testMessageHandlerEvent) {
	l.success = true
}

type testMessageHandlerCommand struct{}

func (e testMessageHandlerCommand) TargetAggregateId() string { return "" }

type testMessageHandlerAggregate struct {
	success bool
}

func (a *testMessageHandlerAggregate) Handle(e testMessageHandlerCommand) ([]Event, error) {
	return []Event{testMessageHandlerEvent{}}, nil
}
func (a *testMessageHandlerAggregate) HandleEvent(e testMessageHandlerEvent) {
	a.success = true
}
