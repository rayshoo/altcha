package analytics

import (
	"context"
	"database/sql"
	"time"
)

type Summary struct {
	Challenges    int64   `json:"challenges"`
	Verified      int64   `json:"verified"`
	Failed        int64   `json:"failed"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
	Errors4XX     int64   `json:"errors_4xx"`
	Errors5XX     int64   `json:"errors_5xx"`
	TotalRequests int64   `json:"total_requests"`
}

type TimeseriesPoint struct {
	Date       string  `json:"date"`
	Challenges int64   `json:"challenges"`
	Verified   int64   `json:"verified"`
	Failed     int64   `json:"failed"`
	AvgLatency float64 `json:"avg_latency"`
}

type LocationEntry struct {
	Continent string  `json:"continent"`
	Country   string  `json:"country"`
	Count     int64   `json:"count"`
	Percent   float64 `json:"percent"`
}

func QuerySummary(ctx context.Context, db *sql.DB, from, to time.Time) (*Summary, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE endpoint = 'challenge') AS challenges,
			COUNT(*) FILTER (WHERE endpoint = 'verify' AND status = 202) AS verified,
			COUNT(*) FILTER (WHERE endpoint = 'verify' AND status = 417) AS failed,
			COALESCE(AVG(latency_ms), 0) AS avg_latency_ms,
			COUNT(*) FILTER (WHERE status >= 400 AND status < 500) AS errors_4xx,
			COUNT(*) FILTER (WHERE status >= 500) AS errors_5xx,
			COUNT(*) AS total_requests
		FROM events
		WHERE timestamp >= $1 AND timestamp < $2
	`
	var s Summary
	err := db.QueryRowContext(ctx, query, from, to).Scan(
		&s.Challenges, &s.Verified, &s.Failed,
		&s.AvgLatencyMs, &s.Errors4XX, &s.Errors5XX, &s.TotalRequests,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func QueryTimeseries(ctx context.Context, db *sql.DB, from, to time.Time) ([]TimeseriesPoint, error) {
	query := `
		SELECT
			DATE(timestamp) AS date,
			COUNT(*) FILTER (WHERE endpoint = 'challenge') AS challenges,
			COUNT(*) FILTER (WHERE endpoint = 'verify' AND status = 202) AS verified,
			COUNT(*) FILTER (WHERE endpoint = 'verify' AND status = 417) AS failed,
			COALESCE(AVG(latency_ms), 0) AS avg_latency
		FROM events
		WHERE timestamp >= $1 AND timestamp < $2
		GROUP BY DATE(timestamp)
		ORDER BY date
	`
	rows, err := db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []TimeseriesPoint
	for rows.Next() {
		var p TimeseriesPoint
		var d time.Time
		if err := rows.Scan(&d, &p.Challenges, &p.Verified, &p.Failed, &p.AvgLatency); err != nil {
			return nil, err
		}
		p.Date = d.Format("2006-01-02")
		points = append(points, p)
	}
	return points, rows.Err()
}

func QueryLocations(ctx context.Context, db *sql.DB, from, to time.Time) ([]LocationEntry, error) {
	query := `
		SELECT
			COALESCE(continent, 'Unknown') AS continent,
			COALESCE(country, 'Unknown') AS country,
			COUNT(*) AS count
		FROM events
		WHERE timestamp >= $1 AND timestamp < $2
		GROUP BY continent, country
		ORDER BY count DESC
	`
	rows, err := db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LocationEntry
	var total int64
	for rows.Next() {
		var e LocationEntry
		if err := rows.Scan(&e.Continent, &e.Country, &e.Count); err != nil {
			return nil, err
		}
		total += e.Count
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range entries {
		if total > 0 {
			entries[i].Percent = float64(entries[i].Count) / float64(total) * 100
		}
	}
	return entries, nil
}
