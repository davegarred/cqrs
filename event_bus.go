package cqrs

import (
	"reflect"
)

type EventBus struct {
	queryEventListeners map[reflect.Type][]*queryEventListener
}

func NewEventBus() *EventBus {
	return &EventBus{make(map[reflect.Type][]*queryEventListener)}
}

func (eventBus *EventBus) RegisterQueryEventHandlers(listener interface{}) {
	aggregateType := reflect.TypeOf(listener)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		f := aggregateType.Method(i)

		if hasEventListenerSignature(f) {
			eventType := f.Type.In(1)
			queryEventListeners := eventBus.queryEventListeners[eventType]
			if queryEventListeners == nil {
				queryEventListeners = make([]*queryEventListener, 0)
			}
			queryEventListeners = append(queryEventListeners, NewEventListener(listener, f))
			eventBus.queryEventListeners[eventType] = queryEventListeners
		}
	}
}

func (eventBus *EventBus) PublishEvents(events []Event) {
	for _, event := range events {
		for _, listener := range eventBus.queryEventListeners[reflect.TypeOf(event)] {
			listener.applyEvent(event)
		}
	}
}
