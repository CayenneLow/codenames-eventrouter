package eventrouter

import (
	"fmt"
	"net"
	"time"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type EventRouter struct {
	config             config.Config
	clientTypeToClient map[client.ClientType]([]client.IClient)
	addrToClientType   map[net.Addr](client.ClientType)
}

func NewEventRouter(config config.Config) EventRouter {
	eventRouter := EventRouter{
		config:             config,
		clientTypeToClient: map[client.ClientType][]client.IClient{},
		addrToClientType:   map[net.Addr]client.ClientType{},
	}
	return eventRouter
}

func (er *EventRouter) AddClient(clientType client.ClientType, cl client.IClient) {
	er.clientTypeToClient[clientType] = append(er.clientTypeToClient[clientType], cl)
	er.addrToClientType[cl.RemoteAddr()] = clientType
	event, err := event.FromJSON([]byte(fmt.Sprintf(`{
		"type": "startConn",
		"gameID": "",
		"timestamp": %d,
		"payload": {
			"status": "success",
			"message": {}
		}
	}`, time.Now().Unix())))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error creating startConn Ack JSON"))
	}
	cl.EmitEvent(event)
}

func (er *EventRouter) RemoveClient(addr net.Addr) error {
	log.Debugf("Removing client: %s", addr)
	if _, ok := er.addrToClientType[addr]; !ok {
		return errors.New(fmt.Sprintf("Client %s does not exist", addr))
	}
	clientType := er.addrToClientType[addr]
	delete(er.addrToClientType, addr)
	clients := er.clientTypeToClient[clientType]
	for i, client := range clients {
		if client.RemoteAddr() == addr {
			// Deletes this client by replacing the current index with the last client in the list
			// then shortening the list by 1
			er.clientTypeToClient[clientType][i] = clients[len(clients)-1]
			er.clientTypeToClient[clientType] = er.clientTypeToClient[clientType][:len(clients)-1]
		}
	}
	return nil
}

func (er *EventRouter) HandleEvent(conn *websocket.Conn, event event.Event) {
	eventType := event.Type
	var recipients []client.IClient
	if eventType == "startConn" {
		var clientType string
		if n, ok := event.Payload.Message["clientType"].(string); ok {
			clientType = string(n)
		}
		cl := client.NewClient(client.GetClientType(clientType), conn, conn.RemoteAddr())
		er.AddClient(client.GetClientType(clientType), cl)
		log.Debugf("Adding %s to clients. Clients: %v", conn.RemoteAddr(), er.clientTypeToClient)
	} else {
		log.Debugf("Received event: %s from client: %v for Game: %s", event.Type, conn.RemoteAddr(), event.GameID)
		if event.Payload.Status == "" {
			// initiator message
			receivers := er.config.GetReceivers(eventType)
			log.Debugf("Receivers: %v", receivers)
			for _, r := range receivers {
				recipients = append(recipients, er.clientTypeToClient[client.GetClientType(r)]...)
			}
		} else {
			// acknowledge mesasge
			acknowledgers := er.config.GetAcknowledgers(eventType)
			log.Debugf("Acknowledgers: %v", acknowledgers)
			for _, a := range acknowledgers {
				recipients = append(recipients, er.clientTypeToClient[client.GetClientType(a)]...)
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
}

// TODO: Remove, router should not be creating new game
// func newGameId() string {
// 	newUuid := uuid.NewString()
// 	gameID := strings.ToUpper(newUuid[:5])
// 	return gameID
// }
