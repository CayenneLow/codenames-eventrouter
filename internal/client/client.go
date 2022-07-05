package client

import (
	"encoding/json"
	"strings"

	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ClientType int

type Client struct {
	Type ClientType
	Ws   *websocket.Conn
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
	switch strings.ToLower(t) {
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
	log.Debugf("Emitting event: %s", eventJson)
	if err := c.Ws.WriteMessage(websocket.TextMessage, eventJson); err != nil {
		return errors.Wrap(err, "Error writing message to websocket")
	}
	return nil
}

func (c *Client) GetType() ClientType {
	return c.Type
}

func (c *Client) GetConn() *websocket.Conn {
	return c.Ws
}
