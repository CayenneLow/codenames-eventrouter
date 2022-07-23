package internal

import (
	"fmt"
	"net/http"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/database"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	eventrouter "github.com/CayenneLow/codenames-eventrouter/internal/event_router"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

type Server struct {
	config      config.Config
	eventRouter eventrouter.EventRouter
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func StartServer(cfg config.Config, db database.Database) {
	log.Info("Starting Server")
	eventRouter := eventrouter.NewEventRouter(cfg, db)
	server := Server{
		config:      cfg,
		eventRouter: eventRouter,
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/ws", server.ws)
	wsEndpoint := fmt.Sprintf(":%s", cfg.WsPort)
	log.Infof("Listening on: %s", wsEndpoint)
	log.Fatal(http.ListenAndServe(wsEndpoint, nil))
}

func (s *Server) ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err, "Unable to upgrade to websocket")
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		log.Debugf("Received message: %s", message)
		if err != nil {
			log.Error(err)
			break
		}
		event, err := event.FromJSON(message)
		if err != nil {
			log.Error(err)
			break
		}
		log.Debugf("Received event: %s from client: %v for Game: %s", event.Type, c.RemoteAddr(), event.GameID)
		s.eventRouter.HandleEvent(c, event)
	}
	// Remove from event router
	s.eventRouter.RemoveClient(c.RemoteAddr())
}
