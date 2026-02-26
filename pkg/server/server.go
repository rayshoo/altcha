package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"altcha/pkg/config"
	"altcha/pkg/handler"
	"altcha/pkg/middleware"
	"altcha/pkg/store"
)

func NewAPIServer(cfg *config.Config, s store.Store) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
		Format: "[API] ${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	if len(cfg.CorsOrigin) > 0 {
		e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
			AllowOrigins: cfg.CorsOrigin,
		}))
	} else {
		e.Use(echomw.CORS())
	}

	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	e.GET("/health", handler.Health())
	e.GET("/challenge", handler.Challenge(cfg))
	e.GET("/verify", handler.Verify(cfg, s))

	return e
}

func NewDemoServer(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	if cfg.IsDebug() {
		e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
			Format: "[DEMO] ${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
		}))
	}
	e.Use(middleware.DemoCSP())

	e.GET("/", handler.DemoPage())
	e.GET("/challenge", handler.DemoChallenge(cfg))
	e.POST("/test", handler.DemoTest(cfg))

	return e
}
