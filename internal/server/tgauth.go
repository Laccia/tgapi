package server

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gitlab.figvam.ru/figvam/tgapi/internal/tgap"
)

func TgAuth(ctx echo.Context) error {
	log.Info().Str("Handle", "auth").Msg("Incoming request")
	msg := &tgap.Code{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&msg)
	if err != nil {
		return err
	}
	tgap.CodeCH <- msg.Code
	return nil
}
