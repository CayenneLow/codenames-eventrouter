package event

type Payload struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Event struct {
	Type      string `json:"type"`
	GameID    string `json:"gameID"`
	Timestamp uint64 `json:"timestamp"`
	Payload   Payload
}
