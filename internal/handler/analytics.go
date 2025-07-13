package handler

import (
	"encoding/json"
	"net/http"

	"app/internal/service"
	"app/internal/views/analyticsview"
	"app/pkg/reqcontext"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	stats, err := h.analyticsService.GetEventStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = analyticsview.DashboardPage(analyticsview.DashboardPageProps{
		Ctx:   ctx,
		Stats: stats,
	}).
		Render(w)
}

func (h *AnalyticsHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.analyticsService.GetEventStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// func GetEventStats() (tracker.EventStats, error) {
// 	ctx := context.Background()
// 	stats := tracker.EventStats{}

// 	// Total events
// 	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_event").Scan(&stats.TotalEvents)
// 	if err != nil {
// 		return stats, err
// 	}

// 	// Unique users
// 	err = db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT user_id) FROM user_event").Scan(&stats.UniqueUsers)
// 	if err != nil {
// 		return stats, err
// 	}

// 	// Top events
// 	rows, err := db.QueryContext(ctx, `
// 		SELECT event_name, COUNT(*) as count
// 		FROM user_event
// 		GROUP BY event_name
// 		ORDER BY count DESC
// 		LIMIT 10
// 	`)
// 	if err != nil {
// 		return stats, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var event tracker.EventCount
// 		err := rows.Scan(&event.EventName, &event.Count)
// 		if err != nil {
// 			return stats, err
// 		}
// 		stats.TopEvents = append(stats.TopEvents, event)
// 	}

// 	// Events over time (last 7 days)
// 	rows, err = db.QueryContext(ctx, `
// 		SELECT DATE(occurred_at) as date, COUNT(*) as count
// 		FROM user_event
// 		WHERE occurred_at >= NOW() - INTERVAL '7 days'
// 		GROUP BY DATE(occurred_at)
// 		ORDER BY date
// 	`)
// 	if err != nil {
// 		return stats, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var point tracker.TimeSeriesPoint
// 		var date time.Time
// 		err := rows.Scan(&date, &point.Count)
// 		if err != nil {
// 			return stats, err
// 		}
// 		point.Date = date.Format("Jan 02")
// 		stats.EventsOverTime = append(stats.EventsOverTime, point)
// 	}

// 	// Events by hour
// 	rows, err = db.QueryContext(ctx, `
// 		SELECT EXTRACT(HOUR FROM occurred_at)::int as hour, COUNT(*) as count
// 		FROM user_event
// 		WHERE occurred_at >= NOW() - INTERVAL '24 hours'
// 		GROUP BY EXTRACT(HOUR FROM occurred_at)
// 		ORDER BY hour
// 	`)
// 	if err != nil {
// 		return stats, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var hourly tracker.HourlyCount
// 		err := rows.Scan(&hourly.Hour, &hourly.Count)
// 		if err != nil {
// 			return stats, err
// 		}
// 		stats.EventsByHour = append(stats.EventsByHour, hourly)
// 	}

// 	// Recent events
// 	rows, err = db.QueryContext(ctx, `
// 		SELECT event_name, user_id, occurred_at, COALESCE(context, '') as context
// 		FROM user_event
// 		ORDER BY occurred_at DESC
// 		LIMIT 20
// 	`)
// 	if err != nil {
// 		return stats, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var event tracker.RecentEvent
// 		err := rows.Scan(&event.EventName, &event.UserID, &event.OccurredAt, &event.Context)
// 		if err != nil {
// 			return stats, err
// 		}
// 		stats.RecentEvents = append(stats.RecentEvents, event)
// 	}

// 	return stats, nil
// }
