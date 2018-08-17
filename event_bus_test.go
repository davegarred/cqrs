package cqrs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventBus_RegisterQueryEventHandlers(t *testing.T) {
	eventBus := NewEventBus()
	eventBus.RegisterQueryEventHandlers(&FooBarEventListener{})

	assert.Equal(t, 4, len(eventBus.queryEventListeners))
}
