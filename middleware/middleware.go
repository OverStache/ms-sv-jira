package middleware

import (
	"ms-sv-jira/helper/logger"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("Access-Control-Allow-Headers", "*")
		c.Request().Header.Set("Access-Control-Allow-Origin", "*")
		c.Request().Header.Set("Access-Control-Allow-Methods", "*")

		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Headers", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "*")
		return next(c)
	}
}

func (m *GoMiddleware) Custom(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("Access-Control-Allow-Headers", "*")
		c.Request().Header.Set("Access-Control-Allow-Origin", "*")
		c.Request().Header.Set("Access-Control-Allow-Methods", "*")
		c.Request().Header.Set("Ngrok-Skip-Browser-Warning", "true")

		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Headers", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "*")
		return next(c)
	}
}

// LOG will handle the LOG middleware
func (m *GoMiddleware) Log(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		l := logger.L
		l.Info("Accepted")

		next(c)

		l.Info("[" + strconv.Itoa(c.Response().Status) + "] " + "[" + c.Request().Method + "] " + c.Request().Host + c.Request().URL.String())

		l.Info("Closing")
		return nil
	}
}

// InitMiddleware intialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
