package internal

import (
	"fmt"
	"net/http"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/database"
	eventrouter "github.com/CayenneLow/codenames-eventrouter/internal/event_router"
	"github.com/CayenneLow/codenames-eventrouter/pkg/event"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

type Server struct {
	config      config.Config
	eventRouter eventrouter.EventRouter
}

func StartServer(cfg config.Config, db database.Database) {
	log.Info("Starting Server")
	eventRouter := eventrouter.NewEventRouter(cfg, db)
	server := Server{
		config:      cfg,
		eventRouter: eventRouter,
	}

	http.HandleFunc("/ws", server.ws)
	wsEndpoint := fmt.Sprintf(":%s", cfg.WsPort)
	log.Infof("Listening on: %s", wsEndpoint)
	log.Fatal(http.ListenAndServe(wsEndpoint, nil))
}

func (s *Server) ws(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Received connection: %s", r.RemoteAddr)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err, "Unable to upgrade to websocket")
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		log.WithFields(log.Fields{
			"message": string(message),
		}).Debug("Received Message")
		if err != nil {
			log.Error(err)
			break
		}
		e, err := event.FromJSON(message)
		if err != nil {
			log.Error(err)
			break
		}
		log.Debugf("Received event: %s from client: %v for Game: %s", e.Type, c.RemoteAddr(), e.GameID)
		s.eventRouter.HandleEvent(c, e)
	}
	// Remove from event router
	s.eventRouter.RemoveClient(c.RemoteAddr())
}
