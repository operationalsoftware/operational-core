package tracker

import "time"

type TrackingEvent struct {
	UserID     int
	EventName  string
	OccurredAt time.Time
	Context    string
	MetaData   map[string]interface{}
}

type EventStats struct {
	TotalEvents    int                 `json:"total_events"`
	UniqueUsers    int                 `json:"unique_users"`
	TopEvents      []EventCount        `json:"top_events"`
	EventsOverTime []TimeSeriesPoint   `json:"events_over_time"`
	UserActivity   []UserActivityPoint `json:"user_activity"`
	EventsByHour   []HourlyCount       `json:"events_by_hour"`
	RecentEvents   []RecentEvent       `json:"recent_events"`
}

type EventCount struct {
	EventName string `json:"event_name"`
	Count     int    `json:"count"`
}

type TimeSeriesPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type UserActivityPoint struct {
	UserID     int `json:"user_id"`
	EventCount int `json:"event_count"`
}

type HourlyCount struct {
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

type RecentEvent struct {
	EventName  string    `json:"event_name"`
	UserID     int       `json:"user_id"`
	OccurredAt time.Time `json:"occurred_at"`
	Context    string    `json:"context"`
}
