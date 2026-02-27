package dashboard

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"altcha/pkg/analytics"
)

func summaryHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		from, to := parseTimeRange(c)
		summary, err := analytics.QuerySummary(c.Request().Context(), db, from, to)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, summary)
	}
}

func timeseriesHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		from, to := parseTimeRange(c)
		points, err := analytics.QueryTimeseries(c.Request().Context(), db, from, to)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, points)
	}
}

func locationsHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		from, to := parseTimeRange(c)
		entries, err := analytics.QueryLocations(c.Request().Context(), db, from, to)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, entries)
	}
}

func parseTimeRange(c echo.Context) (time.Time, time.Time) {
	now := time.Now().UTC()
	to := now
	from := now.AddDate(0, 0, -7)

	if v := c.QueryParam("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = t
		}
	}
	if v := c.QueryParam("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			to = t.AddDate(0, 0, 1) // end of day
		}
	}

	return from, to
}
