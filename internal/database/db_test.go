package database

import (
	"context"
	"testing"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	db  Database
	ctx context.Context
}

func (suite *TestSuite) SetupTest() {
	ctx := context.Background()
	db := Init(ctx, config.Init())
	suite.ctx = ctx
	suite.db = db
}

func (suite *TestSuite) TearDownTest() {
	suite.db.Disconnect(suite.ctx)
}

func (suite *TestSuite) TestGet() {
	suite.T().Run("Read non-existing record", func(t *testing.T) {
		gameId := "INVAL"
		events, err := suite.db.GetEventsByGameId(suite.ctx, gameId)
		assert.NoError(t, err)
		assert.Len(t, events, 0)
	})

	suite.T().Run("Read existing record", func(t *testing.T) {
		gameId := "INTT1"
		events, err := suite.db.GetEventsByGameId(suite.ctx, gameId)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		actualEvent := events[0]
		expected := `{
			"type": "joinGame",
			"GameID": "INTT1",
			"sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa2",
			"timestamp": 1658494936,
			"payload": {
				"status": "",
				"message": {
					"clientType": "server"
				}
			}
		}`
		expectedEvent, err := event.FromJSON([]byte(expected))
		assert.NoError(t, err)
		assert.Equal(t, expectedEvent, actualEvent)
	})

	suite.T().Run("Read multiple records", func(t *testing.T) {
		gameId := "INTT2"
		actualEvents, err := suite.db.GetEventsByGameId(suite.ctx, gameId)
		assert.NoError(t, err)
		assert.Len(t, actualEvents, 2)
		expectedEvents := make([]event.Event, 2, 2)
		expectedEvents[0], err = event.FromJSON([]byte(`
			{
				"GameID": "INTT2",
				"type": "joinGame",
				"sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa2",
				"timestamp": 1658494936,
				"payload": {
					"status": "",
					"message": {
						"clientType": "server"
					}
				}
			}
		`))
		assert.NoError(t, err)
		expectedEvents[1], err = event.FromJSON([]byte(`
			{
				"GameID": "INTT2",
				"type": "joinGame",
				"sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa3",
				"timestamp": 1658494937,
				"payload": {
					"status": "",
					"message": {
						"clientType": "spymaster"
					}
				}
			}
		`))
		assert.NoError(t, err)
		assert.Equal(t, expectedEvents[0], actualEvents[0])
		assert.Equal(t, expectedEvents[1], actualEvents[1])
	})
}

func (suite *TestSuite) TestInsert() {
	expectedEvents := make([]event.Event, 2, 2)
	var err error
	expectedEvents[0], err = event.FromJSON([]byte(`{
		"type": "guess",
		"GameID": "INSRT",
		"sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa2",
		"timestamp": 1658494937,
		"payload": {
			"status": "",
			"message": {
				"boardRow": 1,
				"boardCol": 1
			}
		}
	}`))
	expectedEvents[1], err = event.FromJSON([]byte(`{
		"type": "guess",
		"GameID": "INSRT",
		"sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa3",
		"timestamp": 1658494938,
		"payload": {
			"status": "",
			"message": {
				"boardRow": 1,
				"boardCol": 2
			}
		}
	}`))
	assert.NoError(suite.T(), err)
	suite.T().Run("Insert events", func(t *testing.T) {
		for _, event := range expectedEvents {
			err = suite.db.Insert(suite.ctx, event)
			assert.NoError(t, err)
		}
		actualEvents, err := suite.db.GetEventsByGameId(suite.ctx, "INSRT")
		assert.NoError(t, err)
		assert.Len(t, actualEvents, 2)
		assert.Equal(t, actualEvents[0], expectedEvents[0])
		assert.Equal(t, actualEvents[1], expectedEvents[1])
	})
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
