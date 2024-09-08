package pg

import (
	"context"
	"fmt"
	"os"
	"tgapiV2/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type DB struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Appconfig, log zerolog.Logger) *DB {
	return &DB{
		db: newPg(ctx, cfg, log),
	}
}

func newPg(ctx context.Context,
	cfg *config.Appconfig,
	logger zerolog.Logger) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error while trying connect to DB")
	}

	err = db.Ping(ctx)

	log.Info().Str("comp", "PG").Msg("Trying connect to DB")
	if err != nil {
		log.Fatal().Err(err).Msg("Trying connect to DB")

	}

	log.Info().Str("comp", "PG").Msg("connected")

	Start(db)
	log.Info()
	return db

}

func Start(con *pgxpool.Pool) {
	ctx := context.Background()

	tmp, err := con.Query(ctx, msgTable)
	tmp.Close()
	if err != nil {
		os.Exit(1)
	}
	tmp1, err := con.Query(ctx, histTable)
	tmp1.Close()
	if err != nil {
		os.Exit(1)
	}
}

const (
	msgTable = `CREATE TABLE IF NOT EXISTS tgmsg (
		id SERIAL PRIMARY KEY,
		msg JSON NOT NULL, 
		chat INT NOT NULL);`
	histTable = `CREATE TABLE IF NOT EXISTS tghistory (
		id SERIAL PRIMARY KEY,
		msg TEXT NOT NULL);`
)

func (p *DB) AddMsgPG(ctx context.Context, msg []byte, id int64) error {
	// query := `INSERT INTO tgmsg (msg) VALUES (@msg);`
	args := pgx.NamedArgs{
		"msg":  msg,
		"chat": id,
	}

	tag, err := p.db.Exec(ctx, NewMsg, args)
	if err != nil {

		return err

	}

	fmt.Println("\n", tag)

	return nil
}

const (
	NewMsg = `INSERT INTO tgmsg (msg, chat) VALUES (@msg, @chat);`
)
