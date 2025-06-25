package tracker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tracker struct {
	db *pgxpool.Pool
}

func NewTracker(db *pgxpool.Pool) *Tracker {
	return &Tracker{db: db}
}

func (t *Tracker) TrackEvent(ctx context.Context, event TrackingEvent) error {
	metadataJSON, err := json.Marshal(event.MetaData)
	if err != nil {
		return err
	}

	fmt.Println(event)

	_, err = t.db.Exec(ctx, `
INSERT INTO user_event (
	user_id,
	event_name,
	occurred_at,
	context,
	metadata)
VALUES ($1, $2, $3, $4, $5)
	`, event.UserID, event.EventName, event.OccurredAt, event.Context, metadataJSON)

	return err

}
