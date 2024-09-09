package server

import (
	"encoding/json"
	"tgapiV2/internal/sign"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func TgAuth(ctx echo.Context) error {
	log.Info().Str("Handle", "auth").Msg("Incoming request")
	msg := &sign.Code{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&msg)
	if err != nil {
		return err
	}
	sign.CodeCH <- msg.Code
	return nil
}
