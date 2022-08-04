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
	// Unmarshal events if provided
	var nestedEvents []Event
	if _, ok := event.Payload.Message["events"]; ok {
		log.Debug("HERE")
		if es, ok := event.Payload.Message["events"].([]interface{}); ok {
			log.Debug("HERE2")
			for _, e := range es {
				// This is currently an interface{}, need to marshal into a []byte first
				// before we can convert this into an Event
				eM, err := json.Marshal(e)
				if err != nil {
					return Event{}, nil
				}
				var nestedEvent Event
				if err = json.Unmarshal(eM, &nestedEvent); err != nil {
					return Event{}, nil
				}
				nestedEvents = append(nestedEvents, nestedEvent)
			}
			event.Payload.Message["events"] = nestedEvents
		}
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
