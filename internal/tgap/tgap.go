package tgap

import (
	"context"
	"fmt"
	"log"
	"os"
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

		// _, err = client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		// 	Peer: &tg.InputPeerChannel{
		// 		ChannelID:  -1002065788839,
		// 		AccessHash: user.AccessHash,
		// 	},
		// })
		// if err != nil {
		// 	fmt.Println(err)
		// }

		return gaps.Run(ctx, client.API(), user.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {

				err := parser.DialogsParse(ctx, client)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				logger.Info().Str("comp:", "main").Msg("Application started")
				// fmt.Println(user.AccessHash)
				// kk, err := client.API().MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{Limit: 1, OffsetPeer: &tg.InputPeerChannel{ChannelID: 1195476893}})
				// if err != nil {
				// 	fmt.Println(err)
				// }
				// fmt.Println(kk.String())
				// dialogs, err := json.MarshalIndent(kk, "", "")

				// os.WriteFile("utils/dialogs.json", dialogs, 0644)

				// hist, err := client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{Limit: 100,
				// 	Peer: &tg.InputPeerChannel{
				// 		ChannelID:  1195476893,
				// 		AccessHash: -2798995115427651093,
				// 	},
				// })
				// if err != nil {
				// 	fmt.Println(err)
				// }
				// history, err := json.MarshalIndent(hist, "", "")
				// os.WriteFile("utils/history.json", history, 0644)
			},
		})
	})
}
