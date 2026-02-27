package auth

import (
	"altcha/pkg/config"

	"github.com/labstack/echo/v4"
)

type UserInfo struct {
	Username string
	Roles    []string
	Groups   []string
}

type Provider interface {
	Middleware() echo.MiddlewareFunc
	RegisterRoutes(e *echo.Echo)
}

func NewProvider(cfg *config.Config) Provider {
	switch cfg.AuthProvider {
	case "keycloak":
		return NewOIDCProvider(cfg)
	case "basic":
		return NewBasicProvider(cfg)
	default:
		return NewBasicProvider(cfg)
	}
}

func IsAuthorized(user *UserInfo, cfg *config.Config) bool {
	if len(cfg.AuthAllowedUsers) == 0 && len(cfg.AuthAllowedGroups) == 0 && len(cfg.AuthAllowedRoles) == 0 {
		return true
	}

	for _, u := range cfg.AuthAllowedUsers {
		if u == user.Username {
			return true
		}
	}
	for _, ag := range cfg.AuthAllowedGroups {
		for _, ug := range user.Groups {
			if ag == ug {
				return true
			}
		}
	}
	for _, ar := range cfg.AuthAllowedRoles {
		for _, ur := range user.Roles {
			if ar == ur {
				return true
			}
		}
	}
	return false
}
