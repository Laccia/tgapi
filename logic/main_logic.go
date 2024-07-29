package logic

import (
	"context"
	"fmt"
	"myapp/app/configs"
	"myapp/internal/pgquery"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Logic struct {
	db *pgxpool.Pool
}

func New(configs *configs.Server) *Logic {
	return &Logic{
		db: newPg(configs.PGconn.URL),
	}
}

func newPg(url string) *pgxpool.Pool {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		fmt.Println("lolkek")
		os.Exit(1)
	}

	for {
		err := db.Ping(ctx)
		if err != nil {
			log.Info()
			fmt.Println("lala")
			time.Sleep(1 * time.Second)
		}

		if err == nil {
			break
		}
	}

	pgquery.Start(db)
	log.Info()
	fmt.Println("conn?")
	return db
}

func (l *Logic) PGClose() {
	l.db.Close()
}
