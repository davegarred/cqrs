package components

import (
	"github.com/davegarred/cqrs"
	"reflect"
)

type SynchronousEventBus struct {
	queryEventListeners map[reflect.Type][]*queryEventListener
}

func NewEventBus() *SynchronousEventBus {
	return &SynchronousEventBus{make(map[reflect.Type][]*queryEventListener)}
}

func (eventBus *SynchronousEventBus) RegisterQueryEventHandlers(listener interface{}) {
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

func (eventBus *SynchronousEventBus) PublishEvents(events []cqrs.Event) {
	for _, event := range events {
		for _, listener := range eventBus.queryEventListeners[reflect.TypeOf(event)] {
			listener.applyEvent(event)
		}
	}
}
