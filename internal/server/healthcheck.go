package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type healthCheckResponse struct {
	Result string `json:"result"`
}

func HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, healthCheckResponse{
		Result: "ok",
	})
}
