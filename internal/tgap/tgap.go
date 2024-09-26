package tgap

import (
	"context"
	"fmt"
	"log"
	"tgapiV2/internal/config"
	"tgapiV2/internal/parser"
	"tgapiV2/internal/pg"
	"tgapiV2/internal/sign"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func NewClient(ctx context.Context,
	cfg *config.Appconfig,
	logger zerolog.Logger, db *pg.DB, vt *api.Client) error {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	d := tg.NewUpdateDispatcher()

	gaps := updates.New(updates.Config{
		Handler: d,
	})

	flow := auth.NewFlow(sign.Sign{PhoneNumber: cfg.Phone}, auth.SendCodeOptions{AllowAppHash: true})

	client, err := telegram.ClientFromEnvironment(telegram.Options{UpdateHandler: gaps,
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(gaps.Handle),
		}})
	if err != nil {
		return err
	}
	// Setup message update handlers.
	d.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		logger.Info().Any("asd", update.Message).Msg("message")
		msg := update.Message
		err := parser.MessageParse(ctx, msg, db, cfg.Chats)
		if err != nil {
			logger.Err(err).Str("OnChannel", "message").Msg("err while parse message")
		}
		return nil
	})
	d.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		logger.Info().Any("asd", update.Message).Msg("message")
		msg := update.Message
		err := parser.MessageParse(ctx, msg, db, cfg.Chats)
		if err != nil {
			logger.Err(err).Str("OnChannel", "message").Msg("err while parse message")
		}
		return nil
	})

	return client.Run(ctx, func(ctx context.Context) error {

		// Perform auth if no session is available.
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return errors.Wrap(err, "auth")
		}

		// Fetch user info.
		user, err := client.Self(ctx)
		if err != nil {
			return errors.Wrap(err, "call self")
		}

		return gaps.Run(ctx, client.API(), user.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {

				err := parser.DialogsParse(ctx, client, db, cfg.Chats)
				if err != nil {
					fmt.Println(err)
					logger.Err(err).Str("comp", "main").Msg("Error while download history")
				}

				logger.Info().Str("comp:", "main").Msg("History successfully downloaded")

				logger.Info().Str("comp:", "main").Msg("Application started")

			},
		})
	})
}
