package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
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

	e.GET("/healthcheck", HealthCheck)

	return e
}
