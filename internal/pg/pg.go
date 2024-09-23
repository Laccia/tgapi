package pg

import (
	"context"
	"fmt"
	"os"
	"tgapiV2/internal/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
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

	// Start(db)
	err = migration(ctx, db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Info()
	return db

}

func Start(con *pgxpool.Pool) {

	// m, err := migrate.New(
	// 	"./internal/migrations/*sql",
	// 	"postgres://postgres:postgres@localhost:5432/example?sslmode=disable")
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Fatal()
	// }
	// if err := m.Up(); err != nil {
	// 	fmt.Println(err)
	// 	log.Fatal()
	// }

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

func migration(ctx context.Context, db *pgxpool.Pool) error {
	files, err := os.ReadDir("internal/migrations/")
	if err != nil {
		return fmt.Errorf("read migrations dir error: %s", err)
	}

	migrations := []string{}

	if len(files) < 1 {
		return fmt.Errorf("migrations not found")
	}

	for _, v := range files {
		filename := fmt.Sprintf("%s/%s", "internal/migrations", v.Name())
		content, errRead := os.ReadFile(filename)
		if errRead != nil {
			return fmt.Errorf("failed to read migration file: %s, filename: %s", errRead, filename)
		}

		migrations = append(migrations, string(content))
	}

	fmt.Println(migrations)

	if len(migrations) < 1 {
		return fmt.Errorf("migrations not found")
	}

	for _, m := range migrations {
		tx, errTx := db.Begin(context.Background())
		if errTx != nil {
			return fmt.Errorf("%s fail migrations", errTx)
		}

		rows, err := tx.Exec(ctx, m)
		if err != nil {
			tx.Rollback(ctx)

			return fmt.Errorf("%s fail query migration", err)
		}
		rows.String()
		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("%s fail Commit migration", err)
		}
	}

	return err
}
