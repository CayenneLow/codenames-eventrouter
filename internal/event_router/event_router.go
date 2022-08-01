package eventrouter

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/database"
	"github.com/CayenneLow/codenames-eventrouter/pkg/event"
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
	gameIDToClients      map[string]([]client.IClient) // gameID -> []Clients
	addrToClientMetadata map[net.Addr](ClientMetadata)
	db                   database.Database
}

func NewEventRouter(config config.Config, db database.Database) EventRouter {
	eventRouter := EventRouter{
		config:               config,
		gameIDToClients:      map[string][]client.IClient{},
		addrToClientMetadata: map[net.Addr]ClientMetadata{},
		db:                   db,
	}
	return eventRouter
}

func (er *EventRouter) AddClient(gameID string, clientType client.ClientType, cl client.IClient) {
	log.Debug("Adding Client to Game", log.Fields{
		"gameID":     gameID,
		"clientType": clientType,
		"IP":         cl.RemoteAddr(),
	})
	if _, ok := er.gameIDToClients[gameID]; !ok {
		// Initialize
		er.gameIDToClients[gameID] = []client.IClient{}
	}
	er.gameIDToClients[gameID] = append(er.gameIDToClients[gameID], cl)
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
		clients := er.gameIDToClients[gameID]
		for i, client := range clients {
			// Find client to delete
			if client.RemoteAddr() == addr {
				// Deletes this client by replacing the current index with the last client in the list
				// then shortening the list by 1
				er.gameIDToClients[gameID][i] = clients[len(clients)-1]
				er.gameIDToClients[gameID] = er.gameIDToClients[gameID][:len(clients)-1]
				break
			}
		}
	}
	return nil
}

func (er *EventRouter) HandleEvent(conn *websocket.Conn, event event.Event) {
	er.handleEventRouterEvents(conn, event)
	// Save event to DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := er.db.Insert(ctx, event)
	if err != nil {
		// TODO: Push err to dead letter queue
		log.Error("Error inserting event to db", log.Fields{
			"error": err,
			"event": event,
		})
	}
	// Emit event to consumers
	er.emitEvent(event)
}

func (er *EventRouter) handleEventRouterEvents(conn *websocket.Conn, event event.Event) {
	if event.Payload.Status != "" {
		// Only process non-acknowledgement messages
		return
	}
	eventType := event.Type
	gameID := event.GameID
	switch eventType {
	case "joinGame":
		var clientType string
		if n, ok := event.Payload.Message["clientType"].(string); ok {
			clientType = string(n)
		}
		cl := client.NewClient(client.GetClientType(clientType), conn, conn.RemoteAddr())
		er.AddClient(gameID, client.GetClientType(clientType), cl)
		// TODO: Send turn history snapshot
		log.Debugf("Adding %s to EventRouter Clients. Clients: %v", conn.RemoteAddr(), er.gameIDToClients[gameID])
		ackEvent := er.createAckEvent(event, "success", nil)
		log.Debugf("Emitting joinGame ACK event")
		er.emitEvent(ackEvent)
	}
}

func (er *EventRouter) createAckEvent(event event.Event, status string, messages map[string]interface{}) event.Event {
	event.Payload.Status = status
	event.Payload.Message = map[string](interface{}){} // re-initialize payload
	if messages != nil {
		for k, v := range messages {
			event.Payload.Message[k] = v
		}
	}
	return event
}

func (er *EventRouter) emitEvent(event event.Event) {
	gameID := event.GameID
	recipients := er.gameIDToClients[gameID]
	// Emit to recipients
	for _, r := range recipients {
		log.Debugf("Emitting to: %s", r.RemoteAddr())
		err := r.EmitEvent(event)
		if err != nil {
			log.Error(errors.Wrap(err, fmt.Sprintf("Error emitting event to: %s (%v)", r.CType(), r.RemoteAddr())))
		}
	}
}
