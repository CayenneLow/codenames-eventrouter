package event

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Payload struct {
	Status  string                 `json:"status" bson:"status"`
	Message map[string]interface{} `json:"message" bson:"message"`
}

type Event struct {
	GameID    string  `json:"gameID" bson:"gameID"`
	Type      string  `json:"type" bson:"type"`
	SessionID string  `json:"sessionID" bson:"sessionID"`
	Timestamp uint64  `json:"timestamp" bson:"timestamp"`
	Payload   Payload `json:"payload" bson:"payload"`
}

func FromJSON(j []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(j, &event); err != nil {
		return event, err
	}
	return event, nil
}

func (e Event) JsonString() string {
	j, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		log.Errorf("Error converting event %+v to JSON: %v", e, err)
	}
	return string(j)
}

func (e Event) Bson() []byte {
	b, err := bson.Marshal(e)
	if err != nil {
		log.Errorf("Error converting event %+v to BSON: %v", e, err)
	}
	return b
}
