package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.figvam.ru/figvam/tgapi/internal/request"
)

func HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, request.HealthCheckResponse{
		Result: "ok",
	})
}
