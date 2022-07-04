package eventrouter

import (
	"testing"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/stretchr/testify/assert"
)

type MockServer struct{}
type MockHost struct{}

var emittedServer = false

func (m *MockServer) EmitEvent(event event.Event) {
	emittedServer = true
}

var emittedHost = false

func (m *MockHost) EmitEvent(event event.Event) {
	emittedHost = true
}

func TestForwarding(t *testing.T) {
	cfg := config.Config{
		ForwardingRules: map[string]map[string][]string{
			"newGame": {
				"receivers":     []string{"server"},
				"acknowledgers": []string{"host"},
			},
		},
	}
	eventRouter := NewEventRouter(cfg)
	mockServer := MockServer{}
	eventRouter.AddClient(client.Server, &mockServer)
	mockHost := MockHost{}
	eventRouter.AddClient(client.Host, &mockHost)
	t.Run("Test forward to receiver", func(t *testing.T) {
		mockEvent := event.Event{
			Type:      "newGame",
			GameID:    "test",
			Timestamp: 111111,
			Payload: event.Payload{
				Status:  "",
				Message: "Test",
			},
		}
		eventRouter.HandleEvent(mockEvent)
		assert.True(t, emittedServer)
	})

	t.Run("Test forward to acknowledger", func(t *testing.T) {
		mockAckEvent := event.Event{
			Type:      "newGame",
			GameID:    "test",
			Timestamp: 111111,
			Payload: event.Payload{
				Status:  "Success",
				Message: "Test",
			},
		}
		eventRouter.HandleEvent(mockAckEvent)
		assert.True(t, emittedHost)
	})
}
