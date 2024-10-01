package tgap

import (
	"context"
	"log"

	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"gitlab.figvam.ru/figvam/tgapi/internal/config"
	"gitlab.figvam.ru/figvam/tgapi/internal/parser"
	"gitlab.figvam.ru/figvam/tgapi/internal/pg"
	"gitlab.figvam.ru/figvam/tgapi/internal/sign"
)

func NewClient(ctx context.Context,
	cfg *config.Appconfig,
	logger zerolog.Logger, db *pg.DB, vt *api.Client) error {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	flow := auth.NewFlow(sign.Sign{PhoneNumber: cfg.Phone}, auth.SendCodeOptions{AllowFlashCall: true})

	client, err := telegram.ClientFromEnvironment(telegram.Options{NoUpdates: true})
	if err != nil {
		return err
	}

	return client.Run(ctx, func(ctx context.Context) error {

		// Perform auth if no session is available.
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return errors.Wrap(err, "auth")
		}

		mgs, err := parser.NewHistory(ctx, client, db, cfg.Chats)
		if err != nil {
			logger.Err(err).Str("comp:", "main").Msg("Error while download history")
			return err
		}

		err = mgs.DialogsParse(ctx)
		if err != nil {
			logger.Err(err).Str("comp:", "main").Msg("Error while download history")
			return err
		}

		logger.Info().Str("comp:", "main").Msg("History successfully downloaded")

		logger.Info().Str("comp:", "main").Msg("Application started")

		for {
			time.Sleep(10 * time.Second)
			logger.Info().Str("comp:", "main").Msg("Checking for new messages")
			time.Sleep(1 * time.Second)
			err := mgs.DialogsParse(ctx)
			if err != nil {
				logger.Err(err).Str("comp:", "main").Msg("Error while download new messages")
			}
			logger.Info().Str("comp:", "main").Msg("check complete")

		}

	})

}
