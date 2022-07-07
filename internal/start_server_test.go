package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	cfg      config.Config
	ServerWS *websocket.Conn
	HostWS   *websocket.Conn
	Host2WS  *websocket.Conn
}

var (
	startConnMsg = `
		{
			"type": "startConn",
			"gameID": "",
			"timestamp": %d,
			"payload": {
				"status": "",
				"message": {
					"clientType": "%s",
					"sessionId": "testSession"
				}
			}
		}
	`
)

func (suite *TestSuite) SetupTest() {
	// Init config
	suite.cfg = config.Init()
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%s", suite.cfg.WsHost, suite.cfg.WsPort),
		Path:   "/ws",
	}
	// Start Host websocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Not able to create Host WS: %v", err))
	}
	suite.HostWS = conn
	// Start Host2 websocket
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Not able to create Host WS: %v", err))
	}
	suite.Host2WS = conn
	// Start Server websocket
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Not able to create Server WS: %v", err))
	}
	suite.ServerWS = conn
}

func (suite *TestSuite) TearDownTest() {
	log.Debug("Tearing down after test")
	suite.ServerWS.Close()
	suite.HostWS.Close()
}

// This test creates a websocket for integration testing.
// Feature tested:
// 	- Subscribing client to websocket
// 	- Routing events and receiving acknowledgement events with multiple gameIDs
func (suite *TestSuite) TestSubscribe() {
	// Send Start conn message
	startConnJson := fmt.Sprintf(startConnMsg, time.Now().Unix(), client.Host.String())
	err := suite.HostWS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
	assert.NoError(suite.T(), err)

	var expectedStartConnAckJson event.Event // For assertion
	json.Unmarshal([]byte(`{
		"type": "startConn",
		"gameID": "",
		"timestamp": 111111,
		"payload": {
			"status": "success",
			"message": {}
		}
	}`), &expectedStartConnAckJson)
	// receive start conn acknowledge
	_, msg, err := suite.HostWS.ReadMessage()
	assert.NoError(suite.T(), err)
	actualEvent, err := event.FromJSON(msg)
	assert.NoError(suite.T(), err)
	expectedStartConnAckJson.Timestamp = actualEvent.Timestamp // nullify timestamp comparison
	// Test JSON
	assert.Equal(suite.T(), expectedStartConnAckJson, actualEvent)
}

func (suite *TestSuite) TestNewGame() {
	// Subscribe Host
	startConnJson := fmt.Sprintf(startConnMsg, time.Now().Unix(), client.Host.String())
	err := suite.HostWS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
	assert.NoError(suite.T(), err)
	_, _, err = suite.HostWS.ReadMessage() // Acknowledge
	assert.NoError(suite.T(), err)
	// Subscribe Server
	startConnJson = fmt.Sprintf(startConnMsg, time.Now().Unix(), client.Server.String())
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
	assert.NoError(suite.T(), err)
	_, _, err = suite.ServerWS.ReadMessage() // Acknowledge
	assert.NoError(suite.T(), err)
	// Send New Game message
	newGameJson := fmt.Sprintf(`{
		"type": "newGame",
		"gameID": "",
		"timestamp": %d,
		"payload": {
			"status": "",
			"message": {
				"clientType": "host",
				"sessionId": "testSession"
			}
		}
	}`, time.Now().Unix())
	log.Debug("Writing new game message")
	err = suite.HostWS.WriteMessage(websocket.TextMessage, []byte(newGameJson))
	assert.NoError(suite.T(), err)

	// Assert that server received message
	log.Debug("Asserting Server received message")
	msgType, msg, err := suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	// Convert to Event then compare
	expectedNewGameEvent, err := event.FromJSON([]byte(newGameJson))
	assert.NoError(suite.T(), err)
	actualNewGameEvent, err := event.FromJSON(msg)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedNewGameEvent, actualNewGameEvent)

	// Mocking ACK sent by server
	gameID := "T35T3"
	ts := time.Now().Unix()
	expectedNewGameAckJson := fmt.Sprintf(`{
		"type": "newGame",
		"gameID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "success",
			"message": {}
		}
	}`, gameID, ts)
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(expectedNewGameAckJson))
	assert.NoError(suite.T(), err)
	expectedNewGameAckEvent, err := event.FromJSON([]byte(expectedNewGameAckJson))
	assert.NoError(suite.T(), err)

	// Assert that host received Ack with gameID
	msgType, msg, err = suite.HostWS.ReadMessage()
	actualNewGameAckEvent, err := event.FromJSON(msg)
	// Assert game ID rules
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), "", gameID)
	assert.Len(suite.T(), gameID, 5)
	assert.Equal(suite.T(), strings.ToUpper(gameID), gameID) // assert all upper case
	// Assert rest of JSON
	assert.Equal(suite.T(), expectedNewGameAckEvent, actualNewGameAckEvent)
}

func (suite *TestSuite) TestJoinGameMultipleGameIDs() {
	// Subscribe Host1
	startConnJson := fmt.Sprintf(startConnMsg, time.Now().Unix(), client.Host.String())
	err := suite.HostWS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
	assert.NoError(suite.T(), err)
	_, _, err = suite.HostWS.ReadMessage() // Acknowledge
	assert.NoError(suite.T(), err)
	// Subscribe Host2
	err = suite.Host2WS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
	assert.NoError(suite.T(), err)
	_, _, err = suite.Host2WS.ReadMessage() // Acknowledge
	assert.NoError(suite.T(), err)
	// Subscribe Server
	startConnJson = fmt.Sprintf(startConnMsg, time.Now().Unix(), client.Server.String())
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
	assert.NoError(suite.T(), err)
	_, _, err = suite.ServerWS.ReadMessage() // Acknowledge
	assert.NoError(suite.T(), err)

	// New Game Host1
	newGameJson := fmt.Sprintf(`{
		"type": "newGame",
		"gameID": "",
		"timestamp": %d,
		"payload": {
			"status": "",
			"message": {
				"clientType": "host",
				"sessionId": "testSession"
			}
		}
	}`, time.Now().Unix())
	log.Debug("Writing Host 1 new game message")
	err = suite.HostWS.WriteMessage(websocket.TextMessage, []byte(newGameJson))
	assert.NoError(suite.T(), err)
	// Acknowlege Host1's new game request
	gameID := "HOST1"
	ts := time.Now().Unix()
	expectedNewGameAckJson := fmt.Sprintf(`{
		"type": "newGame",
		"gameID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "success",
			"message": {}
		}
	}`, gameID, ts)
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(expectedNewGameAckJson))
	suite.HostWS.ReadMessage() // ack from host 1

	// New Game Host2
	log.Debug("Writing new game message")
	err = suite.Host2WS.WriteMessage(websocket.TextMessage, []byte(newGameJson))
	assert.NoError(suite.T(), err)

	// Acknowledge Host2's new game request
	gameID2 := "HOST2"
	ts = time.Now().Unix()
	expectedNewGameAckJson2 := fmt.Sprintf(`{
		"type": "newGame",
		"gameID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "success",
			"message": {}
		}
	}`, gameID2, ts)
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(expectedNewGameAckJson2))

	// Acknowledge from Host 2
	_, msg, err := suite.Host2WS.ReadMessage()
	assert.NoError(suite.T(), err)
	ackEvent, err := event.FromJSON(msg)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), gameID, ackEvent.GameID)

	// Asserts that Host 1 did not receive the message intended for host 2
	done := make(chan bool)
	go func() {
		_, msg, err = suite.HostWS.ReadMessage()
		event1, err := event.FromJSON(msg)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), gameID, event1.GameID)
		log.Debugf("Received GameID: %s", event1.GameID)
		done <- true
	}()

	// Unblocks Host 1
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(expectedNewGameAckJson))
	assert.NoError(suite.T(), err)
	<-done
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
