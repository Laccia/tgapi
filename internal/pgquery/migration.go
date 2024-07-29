package pgquery

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Start(con *pgxpool.Pool) {
	ctx := context.Background()

	tmp, err := con.Query(ctx, msgTable)
	tmp.Close()
	if err != nil {
		fmt.Println("lole", err)
		os.Exit(1)
	}
	tmp1, err := con.Query(ctx, histTable)
	tmp1.Close()
	if err != nil {
		fmt.Println("lole", err)
		os.Exit(1)
	}
}

const (
	msgTable = `CREATE TABLE IF NOT EXISTS tgmsg (
		id SERIAL PRIMARY KEY,
		msg JSON NOT NULL);`
	histTable = `CREATE TABLE IF NOT EXISTS tghistory (
		id SERIAL PRIMARY KEY,
		msg TEXT NOT NULL);`
)
