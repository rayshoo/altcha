package middleware

import "github.com/labstack/echo/v4"

func DemoCSP() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Response().Header()
			h.Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self' https://cdn.jsdelivr.net; "+
					"connect-src 'self' https://cdn.jsdelivr.net blob:; "+
					"worker-src 'self' blob:; "+
					"style-src 'self' 'unsafe-inline'")
			return next(c)
		}
	}
}
