package pg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.figvam.ru/figvam/tgapi/internal/config"
)

type RawMSG struct {
	Messages []RawData `json:"Messages"`
}

type RawData struct {
	Messages string           `json:"Message"`
	MsgId    int              `json:"ID"`
	Date     uint64           `json:"Date"`
	ID       map[string]int64 `json:"PeerID"`
}

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

	logger.Info().Str("comp", "PG").Msg("Trying connect to DB")
	if err != nil {
		log.Fatal().Err(err).Msg("Trying connect to DB")

	}

	logger.Info().Str("comp", "PG").Msg("connected")

	// Start(db)
	err = migration(ctx, db)
	if err != nil {
		logger.Fatal().Err(err).Str("comp:", "newPg").Msg("Error while create migration")
	}
	logger.Info().Str("comp", "newPg").Msg("migration success")
	return db

}

func (p *DB) Close() {
	p.db.Close()

}

const (
	getLastId  = `SELECT MAX(msg_id) FROM tghistory WHERE chat_id=@chat;`
	checkExist = `SELECT COUNT(*) AS total_rows FROM tghistory WHERE chat_id=@chat;`
)

func (p *DB) AddHistPG(ctx context.Context, msg []byte) (int, error) {
	ms := RawMSG{}
	err := json.Unmarshal(msg, &ms)
	if err != nil {
		return 0, err
	}
	msgs := ms.Messages
	step, err := p.QueryStream(ctx, msgs)
	if err != nil {
		return 0, err
	}
	return step, nil
}

func (p *DB) CheckExist(ctx context.Context, chat int64) (int64, error) {
	var rows int64
	log.Info().Str("comp:", "CheckExist").Any("Chat: ", chat).Msg("Checking availability rows")
	args := pgx.NamedArgs{
		"chat": chat,
	}
	tx, err := p.db.Begin(ctx)

	if err != nil {
		log.Err(err).Str("comp:", "CheckExist").Msg("Error while query availability rows")
	}

	tag, err := tx.Query(ctx, checkExist, args)

	if err != nil {
		log.Err(err).Str("comp:", "CheckExist").Msg("Error while query availability rows")
	}
	tag.Next()
	err = tag.Scan(&rows)
	if err != nil {
		log.Err(err).Str("comp:", "CheckExist").Msg("Error while query availability rows")

	}
	tag.Close()
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Str("comp:", "CheckExist").Msg("Error while query availability rows")
	}
	return rows, nil
}

func (p *DB) HistoryCheck(ctx context.Context, chat int64) (int, error) {
	var msgId int
	args := pgx.NamedArgs{
		"chat": chat,
	}
	tx, err := p.db.Begin(ctx)
	if err != nil {
		log.Err(err).Str("comp:", "HistoryCheck").Msg("Checking max step")
	}

	tag, err := p.db.Query(ctx, getLastId, args)

	if err != nil {
		log.Err(err).Str("comp:", "HistoryCheck").Msg("Checking max step")
	}
	tag.Next()
	err = tag.Scan(&msgId)
	tag.Close()
	if err != nil {
		log.Err(err).Str("comp:", "HistoryCheck").Msg("Checking max step")
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Str("comp:", "HistoryCheck").Msg("Checking max step")
	}

	log.Info().Str("comp:", "StepCheck").Any("Chat:", chat).Any("Step:", msgId).Msg("Checking max step")
	return msgId, nil
}

func (p *DB) QueryStream(ctx context.Context, msgs []RawData) (int, error) {
	var msgid int
	var Message string
	var id int64
	var date uint64
	var step int

	for i := range msgs {

		tx, err := p.db.Begin(ctx)
		if err != nil {
			log.Err(err).Str("comp:", "QueryStream").Msg("Error while begin query")
		}

		Message = msgs[i].Messages
		msgid = msgs[i].MsgId
		if i == 0 {

			step = msgid
		}
		date = msgs[i].Date
		id = msgs[i].ID["ChannelID"]
		stamp := time.Unix(int64(date), 0)
		if Message != "" && msgid != 0 {

			log.Info().Str("comp:", "QueryStream").Any("message:", msgs[i]).Msg("Pushing new messages")

			args := pgx.NamedArgs{
				"msg":      Message,
				"msg_id":   msgid,
				"msg_date": stamp,
				"chat_id":  id,
			}

			tag, err := tx.Query(ctx, HistoryMsg, args)
			if err != nil {
				log.Err(err).Str("comp:", "QueryStream").Msg("Error while begin exec")
			}
			tag.Close()
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Err(err).Str("comp:", "QueryStream").Msg("Error while commit query")
		}

	}
	return step, nil
}

const (
	NewMsg = `INSERT INTO tgmsg (msg, chat) VALUES (@msg, @chat);`

	HistoryMsg = `INSERT INTO tghistory (msg, msg_id, msg_date, chat_id) VALUES (@msg, @msg_id, @msg_date, @chat_id)`
)

func migration(ctx context.Context, db *pgxpool.Pool) error {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir error: %s", err)
	}

	migrations := []string{}
	newMigrations := []string{}

	if len(files) < 1 {
		return fmt.Errorf("migrations not found")
	}

	presentMigrations, err := techTable(ctx, db)
	if err != nil {
		return fmt.Errorf("fail to techTable: %s", err)
	}

	for _, v := range files {
		_, ok := presentMigrations[v.Name()]
		if ok {
			continue
		}
		filename := fmt.Sprintf("%s/%s", "migrations", v.Name())
		content, errRead := os.ReadFile(filename)
		if errRead != nil {
			return fmt.Errorf("failed to read migration file: %s, filename: %s", errRead, filename)
		}

		migrations = append(migrations, string(content))
		newMigrations = append(newMigrations, v.Name())
	}

	if len(migrations) < 1 && len(presentMigrations) == 0 {
		return fmt.Errorf("migrations not found")
	}

	for _, m := range migrations {
		tx, errTx := db.Begin(context.Background())
		if errTx != nil {
			return fmt.Errorf("%s fail migrations", errTx)
		}

		_, err := tx.Exec(ctx, m)
		if err != nil {
			tx.Rollback(ctx)

			return fmt.Errorf("%s fail query migration", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("%s fail Commit migration", err)
		}
	}
	addToTechTable(ctx, db, newMigrations)

	return err
}

const findMigration = "Select migration From tgapi_tech;"

func techTable(ctx context.Context, db *pgxpool.Pool) (map[string]bool, error) {
	createtable := `
	CREATE TABLE IF NOT EXISTS tgapi_tech (
		migration VARCHAR(255) PRIMARY KEY,
		timestamp TIMESTAMP NOT NULL);
	`

	tx, err := db.Begin(context.Background())
	if err != nil {
		return map[string]bool{}, fmt.Errorf("%s fail migrations", err)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, createtable)
	if err != nil {
		return map[string]bool{}, fmt.Errorf("%s fail migrations", err)
	}

	rows, err := tx.Query(ctx, findMigration)
	if err != nil {
		return map[string]bool{}, fmt.Errorf("%s fail migrations", err)
	}

	defer rows.Close()

	result := make(map[string]bool)
	var tmp string

	for rows.Next() {
		err := rows.Scan(&tmp)
		if err != nil {
			return map[string]bool{}, nil
		}
		result[tmp] = true
	}
	tx.Commit(ctx)
	return result, nil
}

const addMigrationsToTechTable = "INSERT INTO tgapi_tech (migration, timestamp) VALUES (@migration, @timestamp);"

func addToTechTable(ctx context.Context, db *pgxpool.Pool, migrations []string) error {
	pool, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer pool.Rollback(ctx)

	for index := 0; index < len(migrations); index++ {
		args := pgx.NamedArgs{
			"migration": migrations[index],
			"timestamp": time.Now().UTC(),
		}

		rows, err := pool.Query(ctx, addMigrationsToTechTable, args)
		if err != nil {
			return fmt.Errorf("%s fail addToTechTable", err)
		}

		rows.Close()
	}

	err = pool.Commit(ctx)
	return err
}
