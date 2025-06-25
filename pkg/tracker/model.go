package tracker

import "time"

type TrackingEvent struct {
	UserID     int
	EventName  string
	OccurredAt time.Time
	Context    string
	MetaData   map[string]interface{}
}
