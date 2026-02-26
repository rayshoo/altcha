package handler

import (
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
)

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Go      string `json:"go"`
}

var Version = "dev"

func Health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, healthResponse{
			Status:  "ok",
			Version: Version,
			Go:      runtime.Version(),
		})
	}
}
