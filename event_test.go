package event_test

import (
	"testing"

	"github.com/parsilver/event"
	"github.com/stretchr/testify/assert"
)

func TestEvent_Name(t *testing.T) {
	e := event.NewEvent("user.created")
	assert.Equal(t, "user.created", e.Name())
}

func TestEvent_WithArguments(t *testing.T) {
	payload := map[string]interface{}{
		"user_id": 123,
		"email":   "john@example.com",
	}

	e := event.NewEvent("user.created", payload)

	args := e.Arguments()
	assert.Equal(t, payload, args)
}

func TestEvent_StopPropagation(t *testing.T) {
	e := event.NewEvent("user.created")

	assert.False(t, e.IsPropagationStopped())

	e.StopPropagation()

	assert.True(t, e.IsPropagationStopped())
}
