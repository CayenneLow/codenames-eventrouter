package eventrouter

import (
	"fmt"
	"net"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ClientMetadata struct {
	cType   client.ClientType
	gameIDs []string
}

type EventRouter struct {
	config               config.Config
	clientTypeToClient   map[string](map[client.ClientType]([]client.IClient)) // gameID -> clientType -> []Clients
	addrToClientMetadata map[net.Addr](ClientMetadata)
}

func NewEventRouter(config config.Config) EventRouter {
	eventRouter := EventRouter{
		config:               config,
		clientTypeToClient:   map[string](map[client.ClientType][]client.IClient){},
		addrToClientMetadata: map[net.Addr]ClientMetadata{},
	}
	return eventRouter
}

func (er *EventRouter) AddClient(gameID string, clientType client.ClientType, cl client.IClient) {
	if _, ok := er.clientTypeToClient[gameID]; !ok {
		// Initialize
		er.clientTypeToClient[gameID] = map[client.ClientType][]client.IClient{}
		er.clientTypeToClient[gameID][client.Host] = []client.IClient{}
		er.clientTypeToClient[gameID][client.Server] = []client.IClient{}
		er.clientTypeToClient[gameID][client.Spymaster] = []client.IClient{}
	}
	er.clientTypeToClient[gameID][clientType] = append(er.clientTypeToClient[gameID][clientType], cl)
	if v, ok := er.addrToClientMetadata[cl.RemoteAddr()]; ok {
		v.gameIDs = append(v.gameIDs, gameID)
		er.addrToClientMetadata[cl.RemoteAddr()] = v
	} else {
		// create new entry
		er.addrToClientMetadata[cl.RemoteAddr()] = ClientMetadata{
			cType:   clientType,
			gameIDs: []string{gameID},
		}
	}
}

func (er *EventRouter) RemoveClient(addr net.Addr) error {
	log.Debugf("Removing client: %s", addr)
	if _, ok := er.addrToClientMetadata[addr]; !ok {
		return errors.New(fmt.Sprintf("Client %s does not exist", addr))
	}
	clientMetadata := er.addrToClientMetadata[addr]
	delete(er.addrToClientMetadata, addr)
	for _, gameID := range clientMetadata.gameIDs {
		clients := er.clientTypeToClient[gameID][clientMetadata.cType]
		for i, client := range clients {
			// Find client to delete
			if client.RemoteAddr() == addr {
				// Deletes this client by replacing the current index with the last client in the list
				// then shortening the list by 1
				er.clientTypeToClient[gameID][clientMetadata.cType][i] = clients[len(clients)-1]
				er.clientTypeToClient[gameID][clientMetadata.cType] = er.clientTypeToClient[gameID][clientMetadata.cType][:len(clients)-1]
				break
			}
		}
	}
	return nil
}

func (er *EventRouter) HandleEvent(conn *websocket.Conn, event event.Event) {
	log.Debugf("Received event: %s from client: %v for Game: %s", event.Type, conn.RemoteAddr(), event.GameID)
	eventType := event.Type
	gameID := event.GameID
	var recipients []client.IClient
	if eventType == "joinGame" {
		var clientType string
		if n, ok := event.Payload.Message["clientType"].(string); ok {
			clientType = string(n)
		}
		cl := client.NewClient(client.GetClientType(clientType), conn, conn.RemoteAddr())
		er.AddClient(gameID, client.GetClientType(clientType), cl)
		log.Debugf("Adding %s to EventRouter Clients. Clients: %v", conn.RemoteAddr(), er.clientTypeToClient[gameID])
	}
	if event.Payload.Status == "" {
		// initiator message
		receivers := er.config.GetReceivers(eventType)
		log.Debugf("Receivers: %v", receivers)
		for _, r := range receivers {
			recipients = append(recipients, er.clientTypeToClient[gameID][client.GetClientType(r)]...)
		}
	} else {
		// acknowledge mesasge
		acknowledgers := er.config.GetAcknowledgers(eventType)
		log.Debugf("Acknowledgers: %v", acknowledgers)
		for _, a := range acknowledgers {
			recipients = append(recipients, er.clientTypeToClient[gameID][client.GetClientType(a)]...)
		}
	}
	for _, r := range recipients {
		log.Debugf("Emitting to: %s", r.RemoteAddr())
		err := r.EmitEvent(event)
		if err != nil {
			log.Error(errors.Wrap(err, fmt.Sprintf("Error emitting event to: %s (%v)", r.CType(), r.RemoteAddr())))
		}
	}
}

// TODO: Remove, router should not be creating new game
// func newGameId() string {
// 	newUuid := uuid.NewString()
// 	gameID := strings.ToUpper(newUuid[:5])
// 	return gameID
// }
