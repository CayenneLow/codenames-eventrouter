package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromJson(t *testing.T) {
	t.Run("Test conversion", func(t *testing.T) {
		eventJson := `{
				"type": "joinGame",
				"gameID": "D3840",
				"sessionID": "b344328f-e0d9-4692-a4dd-ed0d30e3a565",
				"timestamp": 1659362047,
				"payload": {
					"status": "success",
					"message": {
						"events": [
							{	
								"type": "joinGame",
								"gameID": "D3840",
								"sessionID": "b344328f-e0d9-4692-a4dd-ed0d30e3a566",
								"timestamp": 1659362046,
								"payload": {
									"status": "",
									"message": {
										"clientType": "server"
									}
								}
							}	
						]
					}
				}
			}`

		expected := Event{
			Type:      "joinGame",
			GameID:    "D3840",
			SessionID: "b344328f-e0d9-4692-a4dd-ed0d30e3a565",
			Timestamp: 1659362047,
			Payload: Payload{
				Status: "success",
				Message: map[string]interface{}{
					"events": []Event{
						{
							GameID:    "D3840",
							Type:      "joinGame",
							SessionID: "b344328f-e0d9-4692-a4dd-ed0d30e3a566",
							Timestamp: 1659362046,
							Payload: Payload{
								Status: "",
								Message: map[string]interface{}{
									"clientType": "server",
								},
							},
						},
					},
				},
			},
		}

		actual, err := FromJSON([]byte(eventJson))
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}
