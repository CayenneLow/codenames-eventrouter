package internal

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/client"
	"github.com/CayenneLow/codenames-eventrouter/pkg/event"
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
		"sessionID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "",
			"message": {
				"clientType": "%s"
			}
		}
	}`
	joinGameAckJson = `{
		"type": "joinGame",
		"gameID": "%s",
		"sessionID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "%s",
			"message": {}
		}
	}`
	joinGameSpymasterAckJson = `{
		"type": "joinGame",
		"gameID": "%s",
		"sessionID": "%s",
		"timestamp": %d,
		"payload": {
			"status": "%s",
			"message": {
				"team": "%s"
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
	log.Debug("Server joining game")
	joinGameJson := newJoinGameJson(gameID, client.Server.String(), "test-server-session")
	err := suite.ServerWS.WriteMessage(websocket.TextMessage, joinGameJson)
	assert.NoError(suite.T(), err)

	// Assert that Server receives ACK
	log.Debug("Asserting Server received Server ACK")
	joinGameAckJson := newJoinGameAckJson(gameID, "success", "test-server-session")
	msgType, msg, err := suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)

	// Assert that Server receives joinGame Event
	log.Debug("Asserting Server received Server's joinGame event")
	msgType, msg, err = suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameJson, msg)

	// Host joins game next
	log.Debug("Host joining game")
	joinGameJson = newJoinGameJson(gameID, client.Host.String(), "test-host-session")
	err = suite.HostWS.WriteMessage(websocket.TextMessage, []byte(joinGameJson))
	assert.NoError(suite.T(), err)

	// Assert that Server receives joinGame ACK event
	log.Debug("Asserting Server received Host ACK")
	joinGameAckJson = newJoinGameAckJson(gameID, "success", "test-host-session")
	msgType, msg, err = suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)
	// Assert that Host receives joinGame ACK Event
	log.Debug("Asserting Host received Host ACK")
	msgType, msg, err = suite.HostWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)

	// Assert that Server receives joinGame Event
	log.Debug("Asserting Server received Host's joinGame event")
	msgType, msg, err = suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameJson, msg)
	// Assert that Host receives joinGame Event
	log.Debug("Asserting Host received Host's joinGame event")
	msgType, msg, err = suite.HostWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameJson, msg)

	// Spymaster joins game next
	log.Debug("Spymaster joining game")
	joinGameJson = newJoinGameJson(gameID, client.Spymaster.String(), "test-spymaster-session")
	err = suite.SpymasterWS.WriteMessage(websocket.TextMessage, []byte(joinGameJson))
	assert.NoError(suite.T(), err)

	// Assert that Server receives joinGame ACK event
	log.Debug("Asserting Server received Spymaster ACK")
	joinGameAckJson = newJoinGameAckJson(gameID, "success", "test-spymaster-session")
	msgType, msg, err = suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)
	// Assert that Host receives joinGame ACK Event
	log.Debug("Asserting Host received Spymaster ACK")
	msgType, msg, err = suite.HostWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)
	// Assert that Spymaster receives joinGame ACK Event
	log.Debug("Asserting Spymaster received Spymaster ACK")
	msgType, msg, err = suite.SpymasterWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)

	// Assert that Server receives joinGame Event
	log.Debug("Asserting Server received Spymaster's joinGame event")
	msgType, msg, err = suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameJson, msg)
	// Assert that Host receives joinGame Event
	log.Debug("Asserting Host received Spymaster's joinGame event")
	msgType, msg, err = suite.HostWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameJson, msg)
	// Assert that Spymaster receives joinGame Event
	log.Debug("Asserting Spymaster received Spymaster's joinGame event")
	msgType, msg, err = suite.SpymasterWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameJson, msg)

	// Mock Server joinGame ACK
	team := "RED"
	spymasterJoinGameAckJson := newJoinGameSpymasterAckJson(gameID, "success", "test-spymaster-session", team)
	err = suite.ServerWS.WriteMessage(websocket.TextMessage, []byte(spymasterJoinGameAckJson))
	assert.NoError(suite.T(), err)

	// Assert that Server receives joinGame ACK event with team
	log.Debug("Asserting Server received Spymaster ACK with team")
	joinGameAckJson = newJoinGameSpymasterAckJson(gameID, "success", "test-spymaster-session", team)
	msgType, msg, err = suite.ServerWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)
	// Assert that Host receives joinGame ACK Event with team
	log.Debug("Asserting Host received Spymaster ACK with team")
	msgType, msg, err = suite.HostWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)
	// Assert that Spymaster receives joinGame ACK Event with team
	log.Debug("Asserting Spymaster received Spymaster ACK with team")
	msgType, msg, err = suite.SpymasterWS.ReadMessage()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), websocket.TextMessage, msgType)
	assertJsonEvent(suite, joinGameAckJson, msg)
}

func assertJsonEvent(suite *TestSuite, expectedJson []byte, actualJson []byte) {
	expectedEvent, err := event.FromJSON(expectedJson)
	assert.NoError(suite.T(), err)
	actualEvent, err := event.FromJSON(actualJson)
	assert.NoError(suite.T(), err)
	expectedEvent.Timestamp = actualEvent.Timestamp // nullify timestamp comparison
	assert.Equal(suite.T(), expectedEvent, actualEvent)
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

func newJoinGameJson(gameID string, clientType string, sessionID string) []byte {
	return []byte(fmt.Sprintf(joinGameJson, gameID, sessionID, time.Now().Unix(), clientType))
}

func newJoinGameAckJson(gameID string, status string, sessionID string) []byte {
	return []byte(fmt.Sprintf(joinGameAckJson, gameID, sessionID, time.Now().Unix(), status))
}

func newJoinGameSpymasterAckJson(gameID string, status string, sessionID string, team string) []byte {
	return []byte(fmt.Sprintf(joinGameSpymasterAckJson, gameID, sessionID, time.Now().Unix(), status, team))
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
