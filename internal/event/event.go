package event

import (
	"encoding/json"
)

type Payload struct {
	Status  string                 `json:"status"`
	Message map[string]interface{} `json:"message"`
}

type Event struct {
	Type      string  `json:"type"`
	GameID    string  `json:"gameID"`
	Timestamp uint64  `json:"timestamp"`
	Payload   Payload `json:"payload"`
}

func FromJSON(j []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(j, &event); err != nil {
		return event, err
	}
	return event, nil
}
