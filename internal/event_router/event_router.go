package eventrouter

import (
	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
)

type Client interface {
	EmitEvent(event event.Event)
}

type EventRouter struct {
	config  config.Config
	clients map[client.ClientType](Client)
}

func NewEventRouter(config config.Config) EventRouter {
	eventRouter := EventRouter{
		config:  config,
		clients: map[client.ClientType]Client{},
	}
	return eventRouter
}

func (er *EventRouter) AddClient(clientType client.ClientType, cl Client) {
	er.clients[clientType] = cl
}

func (er *EventRouter) HandleEvent(event event.Event) {
	eventType := event.Type
	var recipients []Client
	if event.Payload.Status == "" {
		// initiator message
		receivers := er.config.GetReceivers(eventType)
		for _, r := range receivers {
			recipients = append(recipients, er.clients[client.GetClientType(r)])
		}
	} else {
		// acknowledge mesasge
		acknowledgers := er.config.GetAcknowledgers(eventType)
		for _, a := range acknowledgers {
			recipients = append(recipients, er.clients[client.GetClientType(a)])
		}
	}
	for _, r := range recipients {
		r.EmitEvent(event)
	}
}
