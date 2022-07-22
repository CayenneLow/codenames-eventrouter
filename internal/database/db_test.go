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
		gameId := "T35T1"
		events, err := suite.db.GetEventsByGameId(suite.ctx, gameId)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		actualEvent := events[0]
		expected := `{
			"type": "joinGame",
			"_id": "T35T1",
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
}

func (suite *TestSuite) TestInsert() {
	json := `{
		"type": "guess",
		"_id": "T35T2",
		"sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa2",
		"timestamp": 1658494937,
		"payload": {
			"status": "",
			"message": {
				"boardRow": 1,
				"boardCol": 1
			}
		}
	}`
	suite.T().Run("Insert event", func(t *testing.T) {
		event, err := event.FromJSON([]byte(json))
		assert.NoError(t, err)
		err = suite.db.Insert(suite.ctx, event)
		assert.NoError(t, err)
		events, err := suite.db.GetEventsByGameId(suite.ctx, "T35T2")
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, events[0], event)
	})
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
