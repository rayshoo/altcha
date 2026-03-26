package handler

import (
	"net/http"
	"time"

	altcha "github.com/altcha-org/altcha-lib-go"
	"github.com/labstack/echo/v4"

	"altcha/pkg/config"
)

func Challenge(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		expires := time.Now().Add(time.Duration(cfg.ExpireMinutes) * time.Minute)

		challenge, err := altcha.CreateChallenge(altcha.ChallengeOptions{
			Algorithm: altcha.Algorithm(cfg.Algorithm),
			HMACKey:   cfg.Secret,
			MaxNumber: int64(cfg.MaxNumber),
			Expires:   &expires,
		})
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, challenge)
	}
}
