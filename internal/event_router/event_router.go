package eventrouter

import (
	"fmt"
	"strings"
	"time"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/google/uuid"
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
	gameId := newGameId()
	event, err := event.FromJSON([]byte(fmt.Sprintf(`{
		"type": "startConn",
		"gameID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "success",
			"message": {}
		}
	}`, gameId, time.Now().Unix())))
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
		if event.Payload.Status == "" {
			// initiator message
			receivers := er.config.GetReceivers(eventType)
			for _, r := range receivers {
				recipients = er.clients[client.GetClientType(r)]
			}
		} else {
			// acknowledge mesasge
			acknowledgers := er.config.GetAcknowledgers(eventType)
			for _, a := range acknowledgers {
				recipients = er.clients[client.GetClientType(a)]
			}
		}
		for _, r := range recipients {
			err := r.EmitEvent(event)
			if err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("Error emitting event to: %s (%v)", r.GetType(), r.GetConn())))
			}
		}
	}
}

// TODO: Remove, router should not be creating new game
func newGameId() string {
	newUuid := uuid.NewString()
	gameID := strings.ToUpper(newUuid[:5])
	return gameID
}
