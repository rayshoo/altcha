package dashboard

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"altcha/pkg/auth"
	"altcha/pkg/config"
)

func NewServer(cfg *config.Config, db *sql.DB) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true

	e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
		Format: "[DASHBOARD] ${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	provider := auth.NewProvider(cfg)
	provider.RegisterRoutes(e)

	api := e.Group("/api", provider.Middleware())
	api.GET("/summary", summaryHandler(db))
	api.GET("/timeseries", timeseriesHandler(db))
	api.GET("/locations", locationsHandler(db))

	e.Static("/", "web/dashboard")

	return e, nil
}
