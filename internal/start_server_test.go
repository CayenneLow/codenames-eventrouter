package internal

import (
	"fmt"
	"net/url"
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
	cfg         config.Config
	ServerWS    *websocket.Conn
	HostWS      *websocket.Conn
	Host2WS     *websocket.Conn
	SpymasterWS *websocket.Conn
}

var (
	joinGameJson = `{
		"type": "joinGame",
		"gameID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "",
			"message": {
				"clientType": "%s",
				"sessionID": "%s"
			}
		}
	}`
	joinGameAckJson = `{
		"type": "joinGame",
		"gameID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "%s",
			"message": {
				"sessionID": "%s"
			}
		}
	}`
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
		log.Fatal(fmt.Sprintf("Not able to create Host 1 WS: %v", err))
	}
	suite.HostWS = conn
	// Start Host2 websocket
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Not able to create Host 2 WS: %v", err))
	}
	suite.Host2WS = conn
	// Start Spymaster websocket
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Not able to create Spymaster WS: %v", err))
	}
	suite.SpymasterWS = conn
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
	suite.Host2WS.Close()
	suite.SpymasterWS.Close()
}

func (suite *TestSuite) TestJoinGame() {
	const gameID = "T35T1" // GameID obtained via direct REST API call Host -> Server
	// Server joins game first
	joinGameJson := newJoinGameJson(gameID, client.Server.String(), "")
	err := suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(joinGameJson))
	assert.NoError(suite.T(), err)

	// Host joins game next
	joinGameJson = newJoinGameJson(gameID, client.Host.String(), "")
	err = suite.HostWS.WriteMessage(websocket.TextMessage, []byte(joinGameJson))
	assert.NoError(suite.T(), err)

	// Assert that Host receives ACK
	log.Debug("Asserting Host received ACK")
	joinGameAckJson := newJoinGameAckJson(gameID, "success", "")
	msgType, msg, err := suite.HostWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	// Convert to Event then compare
	expectedJoinGameAckEvent, err := event.FromJSON([]byte(joinGameAckJson))
	assert.NoError(suite.T(), err)
	actualJoinGameAckEvent, err := event.FromJSON(msg)
	assert.NoError(suite.T(), err)
	expectedJoinGameAckEvent.Timestamp = actualJoinGameAckEvent.Timestamp // nullify timestamp equality check
	assert.Equal(suite.T(), expectedJoinGameAckEvent, actualJoinGameAckEvent)

	// Spymaster joins game next
	joinGameJson = newJoinGameJson(gameID, client.Spymaster.String(), "test-session")
	err = suite.SpymasterWS.WriteMessage(websocket.TextMessage, []byte(joinGameJson))
	assert.NoError(suite.T(), err)

	// Assert Server received join game event
	log.Debug("Asserting Server received ACK")
	msgType, msg, err = suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assert.Equal(suite.T(), joinGameJson, string(msg))

	// Mock Server joinGame ACK
	joinGameAckJson = newJoinGameAckJson(gameID, "success", "test-session")
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(joinGameJson))
	assert.NoError(suite.T(), err)

	// Assert Spymaster receives joinGame ACK
	log.Debug("Asserting Spymaster received ACK")
	msgType, msg, err = suite.SpymasterWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	// Convert to Event then compare
	expectedJoinGameAckEvent, err = event.FromJSON([]byte(joinGameAckJson))
	assert.NoError(suite.T(), err)
	actualJoinGameAckEvent, err = event.FromJSON(msg)
	assert.NoError(suite.T(), err)
	expectedJoinGameAckEvent.Timestamp = actualJoinGameAckEvent.Timestamp // nullify timestamp equality check
	assert.Equal(suite.T(), expectedJoinGameAckEvent, actualJoinGameAckEvent)
}

// func (suite *TestSuite) TestJoinGameMultipleGameIDs() {
// 	// Subscribe Host1
// 	startConnJson := fmt.Sprintf(startConnMsg, time.Now().Unix(), client.Host.String())
// 	err := suite.HostWS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
// 	assert.NoError(suite.T(), err)
// 	_, _, err = suite.HostWS.ReadMessage() // Acknowledge
// 	assert.NoError(suite.T(), err)
// 	// Subscribe Host2
// 	err = suite.Host2WS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
// 	assert.NoError(suite.T(), err)
// 	_, _, err = suite.Host2WS.ReadMessage() // Acknowledge
// 	assert.NoError(suite.T(), err)
// 	// Subscribe Server
// 	startConnJson = fmt.Sprintf(startConnMsg, time.Now().Unix(), client.Server.String())
// 	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(startConnJson))
// 	assert.NoError(suite.T(), err)
// 	_, _, err = suite.ServerWS.ReadMessage() // Acknowledge
// 	assert.NoError(suite.T(), err)

// 	// New Game Host1
// 	newGameJson := fmt.Sprintf(`{
// 		"type": "newGame",
// 		"gameID": "",
// 		"timestamp": %d,
// 		"payload": {
// 			"status": "",
// 			"message": {
// 				"clientType": "host",
// 				"sessionId": "testSession"
// 			}
// 		}
// 	}`, time.Now().Unix())
// 	log.Debug("Writing Host 1 new game message")
// 	err = suite.HostWS.WriteMessage(websocket.TextMessage, []byte(newGameJson))
// 	assert.NoError(suite.T(), err)
// 	// Acknowlege Host1's new game request
// 	gameID := "HOST1"
// 	ts := time.Now().Unix()
// 	expectedNewGameAckJson := fmt.Sprintf(`{
// 		"type": "newGame",
// 		"gameID": "%s",
// 		"timestamp": %d,
// 		"payload": {
// 			"status": "success",
// 			"message": {}
// 		}
// 	}`, gameID, ts)
// 	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(expectedNewGameAckJson))
// 	suite.HostWS.ReadMessage() // ack from host 1

// 	// New Game Host2
// 	log.Debug("Writing new game message")
// 	err = suite.Host2WS.WriteMessage(websocket.TextMessage, []byte(newGameJson))
// 	assert.NoError(suite.T(), err)

// 	// Acknowledge Host2's new game request
// 	gameID2 := "HOST2"
// 	ts = time.Now().Unix()
// 	expectedNewGameAckJson2 := fmt.Sprintf(`{
// 		"type": "newGame",
// 		"gameID": "%s",
// 		"timestamp": %d,
// 		"payload": {
// 			"status": "success",
// 			"message": {}
// 		}
// 	}`, gameID2, ts)
// 	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(expectedNewGameAckJson2))

// 	// Acknowledge from Host 2
// 	_, msg, err := suite.Host2WS.ReadMessage()
// 	assert.NoError(suite.T(), err)
// 	ackEvent, err := event.FromJSON(msg)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), gameID, ackEvent.GameID)

// 	// Asserts that Host 1 did not receive the message intended for host 2
// 	done := make(chan bool)
// 	go func() {
// 		_, msg, err = suite.HostWS.ReadMessage()
// 		event1, err := event.FromJSON(msg)
// 		assert.NoError(suite.T(), err)
// 		assert.Equal(suite.T(), gameID, event1.GameID)
// 		log.Debugf("Received GameID: %s", event1.GameID)
// 		done <- true
// 	}()

// 	// Unblocks Host 1
// 	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(expectedNewGameAckJson))
// 	assert.NoError(suite.T(), err)
// 	<-done
// }

func newJoinGameJson(gameID string, clientType string, sessionID string) string {
	return fmt.Sprintf(joinGameJson, gameID, time.Now().Unix(), clientType, sessionID)
}

func newJoinGameAckJson(gameID string, status string, sessionID string) string {
	return fmt.Sprintf(joinGameAckJson, gameID, time.Now().Unix(), status, sessionID)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
