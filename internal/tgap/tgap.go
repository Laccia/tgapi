package tgap

import (
	"context"

	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"gitlab.figvam.ru/figvam/tgapi/internal/config"
	"gitlab.figvam.ru/figvam/tgapi/internal/parser"
	"gitlab.figvam.ru/figvam/tgapi/internal/pg"
)

type TgClient struct {
	cfg    *config.Appconfig
	logger zerolog.Logger
	db     *pg.DB
	vt     *api.Client
	client *telegram.Client
	flow   auth.Flow
}

func New(cfg *config.Appconfig,
	logger zerolog.Logger, db *pg.DB, vt *api.Client) *TgClient {

	flow := auth.NewFlow(Sign{PhoneNumber: cfg.Phone}, auth.SendCodeOptions{AllowFlashCall: true})
	client := telegram.NewClient(cfg.ID, cfg.Hash, telegram.Options{SessionStorage: &session.FileStorage{Path: cfg.File}})
	return &TgClient{logger: logger, db: db, vt: vt, flow: flow, client: client, cfg: cfg}

}

func (t *TgClient) NewClient(ctx context.Context) error {

	return t.client.Run(ctx, t.clientFunc)

}

func (t *TgClient) clientFunc(ctx context.Context) error {
	// Perform auth if no session is available.
	if err := t.client.Auth().IfNecessary(ctx, t.flow); err != nil {
		return errors.Wrap(err, "auth")
	}

	mgs, err := parser.New(ctx, t.client, t.db, t.cfg.Chats)
	if err != nil {
		t.logger.Fatal().Err(err).Str("comp:", "NewHistory").Msg("Error while download history")
		return err
	}

	err = mgs.DialogsParse(ctx)
	if err != nil {
		t.logger.Fatal().Err(err).Str("comp", "DialogParse").Msg("Error while checking availability rows in base")
		return err
	}

	t.logger.Info().Str("comp:", "main").Msg("History successfully downloaded")

	t.logger.Info().Str("comp:", "main").Msg("Application started")

	for {
		time.Sleep(1 * time.Minute)
		t.logger.Info().Str("comp:", "main").Msg("Checking for new messages")
		err := mgs.DialogsParse(ctx)
		if err != nil {
			t.logger.Err(err).Str("comp:", "main").Msg("Error while downloading new messages")
		}
		t.logger.Info().Str("comp:", "main").Msg("check complete")

	}

}
