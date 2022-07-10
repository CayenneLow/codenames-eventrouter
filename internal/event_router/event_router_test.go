package eventrouter

import (
	"net"
	"testing"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type MockServer struct {
	emitted bool
	ws      *websocket.Conn
	remote  net.Addr
}

func (m MockServer) EmitEvent(event event.Event) error {
	m.emitted = true
	return nil
}
func (m MockServer) CType() client.ClientType { return client.Server }
func (m MockServer) WS() *websocket.Conn      { return &websocket.Conn{} }
func (m MockServer) RemoteAddr() net.Addr     { return m.remote }

type MockHost struct {
	emitted bool
	ws      *websocket.Conn
	remote  net.Addr
}

func (m MockHost) EmitEvent(event event.Event) error {
	m.emitted = true
	return nil
}
func (m MockHost) CType() client.ClientType { return client.Server }
func (m MockHost) WS() *websocket.Conn      { return &websocket.Conn{} }
func (m MockHost) RemoteAddr() net.Addr     { return m.remote }

func TestAddClient(t *testing.T) {
	cfg := config.Init()
	t.Run("Test add client with no conflict", func(t *testing.T) {
		gameID := "T35T1"
		eventRouter := NewEventRouter(cfg)
		mockHost := MockHost{remote: &net.IPAddr{IP: net.IPv4(1, 1, 1, 1)}}
		mockClientMetadata := ClientMetadata{cType: client.Host, gameIDs: []string{gameID}}
		eventRouter.AddClient(gameID, client.Host, &mockHost)

		assert.Len(t, eventRouter.gameIDToClients[gameID], 1)
		assert.Equal(t, &mockHost, eventRouter.gameIDToClients[gameID][0])
		assert.Equal(t, mockClientMetadata, eventRouter.addrToClientMetadata[mockHost.RemoteAddr()])
	})

	t.Run("Test appending client", func(t *testing.T) {
		gameID := "T35T1"
		eventRouter := NewEventRouter(cfg)
		mockHost := MockHost{remote: &net.IPAddr{IP: net.IPv4(1, 1, 1, 1)}}
		mockClientMetadata1 := ClientMetadata{cType: client.Host, gameIDs: []string{gameID}}
		mockHost2 := MockHost{remote: &net.IPAddr{IP: net.IPv4(2, 2, 2, 2)}}
		mockClientMetadata2 := ClientMetadata{cType: client.Host, gameIDs: []string{gameID}}

		eventRouter.AddClient(gameID, client.Host, &mockHost)
		eventRouter.AddClient(gameID, client.Host, &mockHost2)

		assert.Len(t, eventRouter.gameIDToClients[gameID], 2)
		assert.Equal(t, &mockHost, eventRouter.gameIDToClients[gameID][0])
		assert.Equal(t, &mockHost2, eventRouter.gameIDToClients[gameID][1])
		assert.Equal(t, mockClientMetadata1, eventRouter.addrToClientMetadata[mockHost.RemoteAddr()])
		assert.Equal(t, mockClientMetadata2, eventRouter.addrToClientMetadata[mockHost2.RemoteAddr()])
	})

	t.Run("Test adding different types of clients", func(t *testing.T) {
		gameID := "T35T1"
		eventRouter := NewEventRouter(cfg)
		mockHost := MockHost{remote: &net.IPAddr{IP: net.IPv4(1, 1, 1, 1)}}
		mockClientMetadata1 := ClientMetadata{cType: client.Host, gameIDs: []string{gameID}}
		eventRouter.AddClient(gameID, client.Host, &mockHost)

		mockServer := MockServer{remote: &net.IPAddr{IP: net.IPv4(2, 2, 2, 2)}}
		mockClientMetadata2 := ClientMetadata{cType: client.Server, gameIDs: []string{gameID}}
		eventRouter.AddClient(gameID, client.Server, &mockServer)

		assert.Len(t, eventRouter.gameIDToClients[gameID], 2)
		assert.Equal(t, &mockHost, eventRouter.gameIDToClients[gameID][0])
		assert.Equal(t, &mockServer, eventRouter.gameIDToClients[gameID][1])
		assert.Equal(t, mockClientMetadata1, eventRouter.addrToClientMetadata[mockHost.RemoteAddr()])
		assert.Equal(t, mockClientMetadata2, eventRouter.addrToClientMetadata[mockServer.RemoteAddr()])
	})
}

func TestRemoveClient(t *testing.T) {
	cfg := config.Init()
	gameID := "T35T1"
	t.Run("Test removing existing client", func(t *testing.T) {
		eventRouter := NewEventRouter(cfg)
		mockHost := MockHost{remote: &net.IPAddr{IP: net.IPv4(1, 1, 1, 1)}}
		eventRouter.AddClient(gameID, client.Host, &mockHost)
		err := eventRouter.RemoveClient(mockHost.RemoteAddr())
		assert.NoError(t, err)
		assert.Len(t, eventRouter.gameIDToClients[gameID], 0)
		assert.NotContains(t, eventRouter.addrToClientMetadata, mockHost)
	})

	t.Run("Test removing non-existing client", func(t *testing.T) {
		eventRouter := NewEventRouter(cfg)
		remote := &net.IPAddr{IP: net.IPv4(1, 1, 1, 1)}
		err := eventRouter.RemoveClient(remote)
		assert.Error(t, err)
	})
}
