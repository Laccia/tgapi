package server

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

const (
	defaultTimeout = 30 * time.Second
)

type Handler struct {
	logger zerolog.Logger
}

func NewHandler(

	logger zerolog.Logger,

) *Handler {
	return &Handler{

		logger: logger,
	}
}

func (h *Handler) SetRoutes() *echo.Echo {

	e := echo.New()

	e.Use(middleware.Recover())

	e.Use(h.setTimeoutMiddleware)

	e.GET("/healthcheck", HealthCheck)

	e.POST("/tgapi/auth", TgAuth)
	return e
}

func (h *Handler) setTimeoutMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
		defer cancel()
		c.SetRequest(c.Request().Clone(ctx))
		return next(c)
	}
}
