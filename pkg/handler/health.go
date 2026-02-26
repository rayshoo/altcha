package handler

import (
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"

	"altcha/pkg/store"
)

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Go      string `json:"go"`
}

var Version = "dev"

func Health(s store.Store) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := s.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, healthResponse{
				Status:  "unavailable",
				Version: Version,
				Go:      runtime.Version(),
			})
		}

		return c.JSON(http.StatusOK, healthResponse{
			Status:  "ok",
			Version: Version,
			Go:      runtime.Version(),
		})
	}
}
