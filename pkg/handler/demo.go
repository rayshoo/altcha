package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"altcha/pkg/config"
)

func DemoPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File("web/demo/index.html")
	}
}

func DemoChallenge(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		url := fmt.Sprintf("http://localhost:%d/challenge", cfg.Port)

		resp, err := http.Get(url)
		if err != nil {
			return c.NoContent(http.StatusBadGateway)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSONBlob(resp.StatusCode, body)
	}
}

func DemoTest(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := c.FormValue("altcha")
		url := fmt.Sprintf("http://localhost:%d/verify?altcha=%s", cfg.Port, payload)

		resp, err := http.Get(url)
		if err != nil {
			return c.NoContent(http.StatusBadGateway)
		}
		defer resp.Body.Close()

		return c.NoContent(resp.StatusCode)
	}
}
