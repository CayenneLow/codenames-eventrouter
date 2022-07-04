package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Message struct {
}

func TestFromJson(t *testing.T) {
	t.Run("Test initiator", func(t *testing.T) {
		eventJson := `{
			"type": "guess",
			"gameID": "test",
			"timestamp": 111111,
			"payload": {
				"status": "",
				"message": {
					"boardRow": 4,
					"boardCol": 3
				}
			}
		}`

		expected := Event{
			Type:      "guess",
			GameID:    "test",
			Timestamp: 111111,
			Payload: Payload{
				Status: "",
				Message: map[string](interface{}){
					"boardRow": float64(4),
					"boardCol": float64(3),
				},
			},
		}

		actual := FromJSON([]byte(eventJson))

		assert.Equal(t, expected, actual)
	})
}
