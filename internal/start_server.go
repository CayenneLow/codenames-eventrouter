package internal

import (
	"fmt"
	"net/http"

	"github.com/CayenneLow/codenames-eventrouter/config"
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

func StartServer() {
	log.Info("Starting Server")
	cfg := config.Init()
	eventRouter := eventrouter.NewEventRouter(cfg)
	server := Server{
		config:      cfg,
		eventRouter: eventRouter,
	}

	http.HandleFunc("/subscribe", server.subscribe)
	wsEndpoint := fmt.Sprintf("%s:%s", cfg.WsHost, cfg.WsPort)
	log.Infof("Listening on: %s", wsEndpoint)
	log.Fatal(http.ListenAndServe(wsEndpoint, nil))
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err, "Unable to upgrade to websocket")
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
		s.eventRouter.HandleEvent(c, event)
	}
}
