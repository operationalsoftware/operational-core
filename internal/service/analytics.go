package service

import (
	"app/internal/repository"
	"app/pkg/tracker"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsService struct {
	db             *pgxpool.Pool
	authRepository *repository.AuthRepository
}

func NewAnalyticsService(
	db *pgxpool.Pool,
	authRepository *repository.AuthRepository,
) *AnalyticsService {
	return &AnalyticsService{
		db:             db,
		authRepository: authRepository,
	}
}

func (s *AnalyticsService) GetEventStats(
	ctx context.Context,
	// input model.VerifyPasswordLoginInput,
) (tracker.EventStats, error) {
	stats := tracker.EventStats{}

	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_event").Scan(&stats.TotalEvents)
	if err != nil {
		return stats, err
	}

	// Unique users
	err = s.db.QueryRow(ctx, "SELECT COUNT(DISTINCT user_id) FROM user_event").Scan(&stats.UniqueUsers)
	if err != nil {
		return stats, err
	}

	// Top events
	rows, err := s.db.Query(ctx, `
		SELECT event_name, COUNT(*) as count 
		FROM user_event 
		GROUP BY event_name 
		ORDER BY count DESC 
		LIMIT 10
	`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var event tracker.EventCount
		err := rows.Scan(&event.EventName, &event.Count)
		if err != nil {
			return stats, err
		}
		stats.TopEvents = append(stats.TopEvents, event)
	}

	// Events over time (last 7 days)
	rows, err = s.db.Query(ctx, `
		SELECT DATE(occurred_at) as date, COUNT(*) as count
		FROM user_event
		WHERE occurred_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(occurred_at)
		ORDER BY date
	`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var point tracker.TimeSeriesPoint
		var date time.Time
		err := rows.Scan(&date, &point.Count)
		if err != nil {
			return stats, err
		}
		point.Date = date.Format("Jan 02")
		stats.EventsOverTime = append(stats.EventsOverTime, point)
	}

	// Events by hour
	rows, err = s.db.Query(ctx, `
		SELECT EXTRACT(HOUR FROM occurred_at)::int as hour, COUNT(*) as count
		FROM user_event
		WHERE occurred_at >= NOW() - INTERVAL '24 hours'
		GROUP BY EXTRACT(HOUR FROM occurred_at)
		ORDER BY hour
	`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var hourly tracker.HourlyCount
		err := rows.Scan(&hourly.Hour, &hourly.Count)
		if err != nil {
			return stats, err
		}
		stats.EventsByHour = append(stats.EventsByHour, hourly)
	}

	// Recent events
	rows, err = s.db.Query(ctx, `
		SELECT event_name, user_id, occurred_at, COALESCE(context, '') as context
		FROM user_event
		ORDER BY occurred_at DESC
		LIMIT 20
	`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var event tracker.RecentEvent
		err := rows.Scan(&event.EventName, &event.UserID, &event.OccurredAt, &event.Context)
		if err != nil {
			return stats, err
		}
		stats.RecentEvents = append(stats.RecentEvents, event)
	}

	return stats, nil
}
