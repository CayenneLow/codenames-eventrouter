package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/CayenneLow/codenames-eventrouter/pkg/event"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ClientType int

type IClient interface {
	EmitEvent(event event.Event) error
	CType() ClientType
	WS() *websocket.Conn
	RemoteAddr() net.Addr
}

type Client struct {
	cType      ClientType
	ws         *websocket.Conn
	remoteAddr net.Addr
}

const (
	Host ClientType = iota
	Spymaster
	Server
	Unknown
)

func NewClient(cType ClientType, ws *websocket.Conn, remoteAddr net.Addr) Client {
	return Client{
		cType:      cType,
		ws:         ws,
		remoteAddr: remoteAddr,
	}
}

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

func (c Client) EmitEvent(event event.Event) error {
	eventJson, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "Error marshalling event")
	}
	log.WithFields(log.Fields{
		"events": string(eventJson),
	}).Debug("Emitting event")
	if err := c.ws.WriteMessage(websocket.TextMessage, eventJson); err != nil {
		return errors.Wrap(err, "Error writing message to websocket")
	}
	return nil
}

func (c Client) CType() ClientType {
	return c.cType
}

func (c Client) WS() *websocket.Conn {
	return c.ws
}

func (c Client) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c Client) String() string {
	return fmt.Sprintf("Client %s (%s)", c.cType.String(), c.ws.RemoteAddr())
}
