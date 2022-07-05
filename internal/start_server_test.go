package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationSubscribe(t *testing.T) {
	// Start server
	cfg := config.Init()
	// Start client websocket
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%s", cfg.WsHost, cfg.WsPort),
		Path:   "/ws",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)

	var startConnAckJson event.Event
	json.Unmarshal([]byte(`{
		"type": "startConn",
		"gameID": "",
		"timestamp": 111111,
		"payload": {
			"status": "success",
			"message": {}
		}
	}`), &startConnAckJson)
	done := make(chan bool)
	// Goroutine to receive message
	go func() {
		index := 0
		defer close(done)
		for {
			switch index {
			case 0:
				// 0: receive start conn acknowledge
				_, msg, err := conn.ReadMessage()
				assert.NoError(t, err)
				actualEvent, err := event.FromJSON(msg)
				assert.NoError(t, err)
				startConnAckJson.Timestamp = actualEvent.Timestamp // nullify timestamp comparison
				assert.Equal(t, startConnAckJson, actualEvent)
				done <- true
				index++

			}
		}
	}()

	// 0: Send Start conn message
	startConnJson := `{
		"type": "startConn",
		"gameID": "",
		"timestamp": 111111,
		"payload": {
			"status": "",
			"message": {
				"clientType": "host",
				"sessionId": "testSession"
			}
		}
	}`
	err = conn.WriteMessage(websocket.TextMessage, []byte(startConnJson))
	assert.NoError(t, err)
	log.Debug("Closing websocket")
	<-done
	conn.Close()
}
