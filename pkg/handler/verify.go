package handler

import (
	"net/http"

	altcha "github.com/altcha-org/altcha-lib-go"
	"github.com/labstack/echo/v4"

	"altcha/pkg/config"
	"altcha/pkg/store"
)

func Verify(cfg *config.Config, s store.Store) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := c.QueryParam("altcha")

		exists, err := s.Exists(payload)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if exists {
			return c.NoContent(http.StatusExpectationFailed)
		}

		ok, err := altcha.VerifySolution(payload, cfg.Secret, true)
		if err != nil {
			return c.NoContent(http.StatusExpectationFailed)
		}

		_ = s.Add(payload)

		if ok {
			return c.NoContent(http.StatusAccepted)
		}
		return c.NoContent(http.StatusExpectationFailed)
	}
}
