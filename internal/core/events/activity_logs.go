package events

import "encoding/json"

type LogActivityEvent struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Action    string          `json:"action"`
	Timestamp string          `json:"timestamp"`
	Metadata  json.RawMessage `json:"metadata"`
}
