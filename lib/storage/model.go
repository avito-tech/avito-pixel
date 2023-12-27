package storage

import "time"

type HitTable struct {
	SessionID string    `db:"sessionID"`
	EventType string    `db:"eventType"`
	EventTime time.Time `json:"event_time" db:"event_time"`
	Platform  string    `db:"platform"`
	Meta      any       `db:"meta"`
}
