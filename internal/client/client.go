package client

import (
	"encoding/json"

	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type ClientType int

type Client struct {
	ws websocket.Conn
}

const (
	Host ClientType = iota
	Spymaster
	Server
	Unknown
)

func (ct ClientType) String() string {
	switch ct {
	case Host:
		return "host"
	case Spymaster:
		return "spymaster"
	case Server:
		return "server"
	default:
		return "unknown"
	}
}

func GetClientType(t string) ClientType {
	switch t {
	case "host":
		return Host
	case "spymaster":
		return Spymaster
	case "server":
		return Server
	default:
		return Unknown
	}
}

func (c *Client) EmitEvent(event event.Event) error {
	eventJson, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "Error marshalling event")
	}
	if err := c.ws.WriteMessage(websocket.TextMessage, eventJson); err != nil {
		return errors.Wrap(err, "Error writing message to websocket")
	}
	return nil
}
