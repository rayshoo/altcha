package handler

import (
	"net/http"
	"sync"

	altcha "github.com/altcha-org/altcha-lib-go"
	"github.com/labstack/echo/v4"

	"altcha/pkg/config"
)

func Verify(cfg *config.Config) echo.HandlerFunc {
	var (
		mu          sync.Mutex
		recordCache []string
	)

	return func(c echo.Context) error {
		payload := c.QueryParam("altcha")

		mu.Lock()
		defer mu.Unlock()

		for _, r := range recordCache {
			if r == payload {
				return c.NoContent(http.StatusExpectationFailed)
			}
		}

		ok, err := altcha.VerifySolution(payload, cfg.Secret, true)
		if err != nil {
			return c.NoContent(http.StatusExpectationFailed)
		}

		recordCache = append(recordCache, payload)
		if len(recordCache) > cfg.MaxRecords {
			recordCache = recordCache[1:]
		}

		if ok {
			return c.NoContent(http.StatusAccepted)
		}
		return c.NoContent(http.StatusExpectationFailed)
	}
}
