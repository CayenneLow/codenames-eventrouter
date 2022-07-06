package eventrouter

import (
	"fmt"
	"time"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	EmitEvent(event event.Event) error
	GetType() client.ClientType
	GetConn() *websocket.Conn
}

type EventRouter struct {
	config  config.Config
	clients map[client.ClientType]([]Client)
}

func NewEventRouter(config config.Config) EventRouter {
	eventRouter := EventRouter{
		config:  config,
		clients: map[client.ClientType][]Client{},
	}
	return eventRouter
}

func (er *EventRouter) AddClient(clientType client.ClientType, cl Client) {
	er.clients[clientType] = append(er.clients[clientType], cl)
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

func (er *EventRouter) HandleEvent(conn *websocket.Conn, event event.Event) {
	eventType := event.Type
	var recipients []Client
	if eventType == "startConn" {
		var clientType string
		if n, ok := event.Payload.Message["clientType"].(string); ok {
			clientType = string(n)
		}
		cl := &client.Client{
			Ws: conn,
		}
		er.AddClient(client.GetClientType(clientType), cl)
		log.Debugf("Adding %s to clients. Clients: %v", conn.RemoteAddr(), er.clients)
	} else {
		log.Debugf("Received event: %s from client: %v for Game: %s", event.Type, conn.RemoteAddr(), event.GameID)
		if event.Payload.Status == "" {
			// initiator message
			receivers := er.config.GetReceivers(eventType)
			log.Debugf("Receivers: %v", receivers)
			for _, r := range receivers {
				recipients = append(recipients, er.clients[client.GetClientType(r)]...)
			}
		} else {
			// acknowledge mesasge
			acknowledgers := er.config.GetAcknowledgers(eventType)
			log.Debugf("Acknowledgers: %v", acknowledgers)
			for _, a := range acknowledgers {
				recipients = append(recipients, er.clients[client.GetClientType(a)]...)
			}
		}
		for _, r := range recipients {
			log.Debugf("Emitting to: %s", r.GetConn().RemoteAddr())
			err := r.EmitEvent(event)
			if err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("Error emitting event to: %s (%v)", r.GetType(), r.GetConn())))
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
