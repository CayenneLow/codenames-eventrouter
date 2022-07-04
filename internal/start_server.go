package eventrouter

import (
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
	cfg := config.Init()
	eventRouter := eventrouter.NewEventRouter(cfg)
	server := Server{
		config:      cfg,
		eventRouter: eventRouter,
	}

	http.HandleFunc("/subscribe", server.subscribe)
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err, "Unable to upgrade to websocket")
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error(err)
			break
		}
		event := event.FromJSON(message)
		s.eventRouter.HandleEvent(c, event)
	}
}
