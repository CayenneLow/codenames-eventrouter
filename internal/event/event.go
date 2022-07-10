package event

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Payload struct {
	Status  string                 `json:"status"`
	Message map[string]interface{} `json:"message"`
}

type Event struct {
	Type      string  `json:"type"`
	GameID    string  `json:"gameID"`
	SessionID string  `json:"sessionID"`
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

func (e Event) String() string {
	j, err := json.MarshalIndent(e, "", "	")
	if err == nil {
		log.Errorf("Error converting event %v to String: %v", e, err)
	}
	return string(j)
}
