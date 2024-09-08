package server

import (
	"net/http"
	"tgapiV2/internal/request"

	"github.com/labstack/echo/v4"
)

func HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, request.HealthCheckResponse{
		Result: "ok",
	})
}
