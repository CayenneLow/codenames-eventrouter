package eventrouter

import (
	"testing"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type MockServer struct{}
type MockHost struct{}

var emittedServer = false

func (m *MockServer) EmitEvent(event event.Event) error {
	emittedServer = true
	return nil
}
func (m *MockServer) GetType() client.ClientType { return client.Server }
func (m *MockServer) GetConn() *websocket.Conn   { return &websocket.Conn{} }

var emittedHost = false

func (m *MockHost) EmitEvent(event event.Event) error {
	emittedHost = true
	return nil
}
func (m *MockHost) GetType() client.ClientType { return client.Host }
func (m *MockHost) GetConn() *websocket.Conn   { return &websocket.Conn{} }

func TestForwarding(t *testing.T) {
	cfg := initConfig()
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
				Status: "",
				Message: map[string](interface{}){
					"Test": "Test",
				},
			},
		}
		eventRouter.HandleEvent(&websocket.Conn{}, mockEvent)
		assert.True(t, emittedServer)
	})

	t.Run("Test forward to acknowledger", func(t *testing.T) {
		mockAckEvent := event.Event{
			Type:      "newGame",
			GameID:    "test",
			Timestamp: 111111,
			Payload: event.Payload{
				Status: "Success",
				Message: map[string](interface{}){
					"Test": "Test",
				},
			},
		}
		eventRouter.HandleEvent(&websocket.Conn{}, mockAckEvent)
		assert.True(t, emittedHost)
	})
}

func TestAddClient(t *testing.T) {
	cfg := initConfig()
	t.Run("Test add client with no conflict", func(t *testing.T) {
		eventRouter := NewEventRouter(cfg)
		mockHost := MockHost{}
		eventRouter.AddClient(client.Host, &mockHost)

		assert.Len(t, eventRouter.clients[client.Host], 1)
		assert.Equal(t, &mockHost, eventRouter.clients[client.Host][0])
	})

	t.Run("Test appending client", func(t *testing.T) {
		eventRouter := NewEventRouter(cfg)
		mockHost := MockHost{}
		eventRouter.AddClient(client.Host, &mockHost)
		eventRouter.AddClient(client.Host, &mockHost)

		assert.Len(t, eventRouter.clients[client.Host], 2)
		assert.Equal(t, &mockHost, eventRouter.clients[client.Host][0])
		assert.Equal(t, &mockHost, eventRouter.clients[client.Host][1])
	})

	t.Run("Test adding different types of clients", func(t *testing.T) {
		eventRouter := NewEventRouter(cfg)
		mockHost := MockHost{}
		eventRouter.AddClient(client.Host, &mockHost)
		mockServer := MockServer{}
		eventRouter.AddClient(client.Server, &mockServer)

		assert.Len(t, eventRouter.clients[client.Host], 1)
		assert.Len(t, eventRouter.clients[client.Server], 1)
		assert.Equal(t, &mockHost, eventRouter.clients[client.Host][0])
		assert.Equal(t, &mockServer, eventRouter.clients[client.Server][0])
	})
}

func initConfig() config.Config {
	return config.Config{
		ForwardingRules: map[string]map[string][]string{
			"newGame": {
				"receivers":     []string{"server"},
				"acknowledgers": []string{"host"},
			},
		},
	}
}
