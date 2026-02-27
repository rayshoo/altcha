package analytics

import (
	"time"

	"github.com/labstack/echo/v4"
)

func Middleware(collector *Collector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			if path != "/challenge" && path != "/verify" {
				return next(c)
			}

			start := time.Now()
			err := next(c)
			latency := time.Since(start).Seconds() * 1000

			collector.Record(Event{
				Timestamp: start,
				Endpoint:  path[1:], // strip leading /
				ClientIP:  c.RealIP(),
				Status:    c.Response().Status,
				LatencyMs: latency,
			})

			return err
		}
	}
}
