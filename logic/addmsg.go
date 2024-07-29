package logic

import (
	"context"
	"fmt"
	"myapp/logic/pgq"
	"os"

	"github.com/jackc/pgx/v5"
)

func (l *Logic) AddMsg(ctx context.Context, msg []byte) error {
	// query := `INSERT INTO tgmsg (msg) VALUES (@msg);`
	args := pgx.NamedArgs{
		"msg": msg,
	}

	tag, err := l.db.Exec(ctx, pgq.AddMsg, args)
	if err != nil {

		fmt.Println(tag)
		fmt.Println(err)
		os.Exit(1)
	}
	// tag.Close()
	fmt.Println("ZDESS\n", tag)

	return nil
}

func (l *Logic) AddHistory(ctx context.Context, msg string) error {
	// query := `INSERT INTO tgmsg (msg) VALUES (@msg);`
	args := pgx.NamedArgs{
		"msg": msg,
	}

	tag, err := l.db.Exec(ctx, pgq.AddHistory, args)
	if err != nil {

		fmt.Println(tag)
		fmt.Println(err)
		os.Exit(1)
	}
	// tag.Close()
	fmt.Println("ZDESS\n", tag)

	return nil
}
