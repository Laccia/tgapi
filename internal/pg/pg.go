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

const (
	getLastId  = `SELECT MAX(msg_id) FROM tghistory WHERE chat_id=@chat;`
	checkExist = `SELECT COUNT(*) AS total_rows FROM tghistory Where chat_id=@chat; `
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

func (p *DB) AddAllHistPG(ctx context.Context, msg []byte) (int, error) {
	fmt.Println("ok")
	ms := RawMSG{}

	err := json.Unmarshal(msg, &ms)
	if err != nil {
		return 0, err
	}
	msgs := ms.Messages
	step, err := p.AllQueryStream(ctx, msgs)
	if err != nil {
		return 0, err
	}

	return step, nil
}

func (p *DB) AddHistPG(ctx context.Context, msg []byte) (int, error) {
	fmt.Println("ok")
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
	fmt.Println(chat)
	args := pgx.NamedArgs{
		"chat": chat,
	}

	tag, err := p.db.Query(ctx, checkExist, args)
	if err != nil {
		return 0, err
	}
	tag.Next()
	err = tag.Scan(&rows)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return rows, nil
}

func (p *DB) HistoryCheck(ctx context.Context, chat int64) (int, error) {
	var msgId int
	fmt.Println(chat)
	args := pgx.NamedArgs{
		"chat": chat,
	}

	tag, err := p.db.Query(ctx, getLastId, args)
	if err != nil {
		return 0, err
	}
	tag.Next()
	err = tag.Scan(&msgId)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return msgId, nil
}

func (p *DB) AllQueryStream(ctx context.Context, msgs []RawData) (int, error) {
	var Msgid int
	var Message string
	var id int64
	var date uint64
	for i := range msgs {
		fmt.Println(i)
		Message = msgs[i].Messages
		Msgid = msgs[i].MsgId
		date = msgs[i].Date
		id = msgs[i].ID["ChannelID"]
		fmt.Println(Msgid)
		stamp := time.Unix(int64(date), 0)
		if Message != "" && Msgid != 0 {
			args := pgx.NamedArgs{
				"msg":      Message,
				"msg_id":   Msgid,
				"msg_date": stamp,
				"chat_id":  id,
			}

			tag, err := p.db.Exec(ctx, HistoryMsg, args)
			if err != nil {
				return 0, err
			}
			fmt.Println("\n", tag)
		}

	}
	fmt.Println("pg", Msgid)
	return Msgid, nil
}

func (p *DB) QueryStream(ctx context.Context, msgs []RawData) (int, error) {
	var msgid int
	var Message string
	var id int64
	var date uint64
	var step int
	for i := range msgs {
		fmt.Println(i)
		Message = msgs[i].Messages
		msgid = msgs[i].MsgId
		if i == 0 {
			step = msgid
		}
		date = msgs[i].Date
		id = msgs[i].ID["ChannelID"]
		fmt.Println(msgid)
		stamp := time.Unix(int64(date), 0)
		if Message != "" && msgid != 0 {
			args := pgx.NamedArgs{
				"msg":      Message,
				"msg_id":   msgid,
				"msg_date": stamp,
				"chat_id":  id,
			}

			tag, err := p.db.Exec(ctx, HistoryMsg, args)
			if err != nil {
				return 0, err
			}
			fmt.Println("\n", tag)
		}

	}
	fmt.Println("pg", step)
	return step, nil
}

const (
	NewMsg = `INSERT INTO tgmsg (msg, chat) VALUES (@msg, @chat);`

	HistoryMsg = `INSERT INTO tghistory (msg, msg_id, msg_date, chat_id) VALUES (@msg, @msg_id, @msg_date, @chat_id)`
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

		rows, err := tx.Query(ctx, m)
		if err != nil {
			tx.Rollback(ctx)

			return fmt.Errorf("%s fail query migration", err)
		}
		rows.Close()
		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("%s fail Commit migration", err)
		}
	}

	return err
}
