package auth

import (
	"crypto/subtle"
	"net/http"

	"altcha/pkg/config"

	"github.com/labstack/echo/v4"
)

type BasicProvider struct {
	cfg *config.Config
}

func NewBasicProvider(cfg *config.Config) *BasicProvider {
	return &BasicProvider{cfg: cfg}
}

func (p *BasicProvider) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, pass, ok := c.Request().BasicAuth()
			if !ok {
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="ALTCHA Dashboard"`)
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}

			userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(p.cfg.AuthUsername)) == 1
			passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(p.cfg.AuthPassword)) == 1

			if !userMatch || !passMatch {
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="ALTCHA Dashboard"`)
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}

			userInfo := &UserInfo{Username: user}
			if !IsAuthorized(userInfo, p.cfg) {
				return c.String(http.StatusForbidden, "Forbidden")
			}

			c.Set("user", userInfo)
			return next(c)
		}
	}
}

func (p *BasicProvider) RegisterRoutes(e *echo.Echo) {
	e.GET("/auth/logout", func(c echo.Context) error {
		c.Response().Header().Set("WWW-Authenticate", `Basic realm="ALTCHA Dashboard"`)
		return c.String(http.StatusUnauthorized, "Logged out. Close this tab or re-enter credentials.")
	})
}
