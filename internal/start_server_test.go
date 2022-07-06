package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationSubscribe(t *testing.T) {
	// Init config
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
				// test gameID
				gameID := actualEvent.GameID
				assert.NotEqual(t, "", gameID)
				assert.Len(t, gameID, 5)
				assert.Equal(t, strings.ToUpper(gameID), gameID) // assert all upper case
				startConnAckJson.GameID = gameID                 // nullify timestamp comparison
				// Test JSON
				assert.Equal(t, startConnAckJson, actualEvent)
				// Mark this test as done
				index++
			default:
				done <- true
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
