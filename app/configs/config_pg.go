package configs

import (
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog/log"
)

type PGconn struct {
	URL    string `json:"url"`
	pgName string
	pgPass string
	pgAddr string
	pgPort string
	PgBase string `json:"base"`
}

func NewPG() *PGconn {
	result := PGconn{}

	var tmp bool

	result.pgName, tmp = os.LookupEnv("PG_Name")
	if !tmp {
		log.Error()
		os.Exit(1)
	}

	result.pgPass, tmp = os.LookupEnv("PG_Pass")
	if !tmp {
		log.Error()
		os.Exit(1)
	}

	result.pgAddr, tmp = os.LookupEnv("PG_Add")
	if !tmp {
		log.Error()
		os.Exit(1)
	}

	result.PgBase, tmp = os.LookupEnv("PG_Base")
	if !tmp {
		log.Error()
		os.Exit(1)
	}

	result.pgPort, tmp = os.LookupEnv("PG_Port")
	if !tmp {
		log.Info()

		result.pgPort = "5432"
	}

	result.URL = fmt.Sprintf("postgres://%s:%s@%s/%s",
		result.pgName,
		result.pgPass,
		net.JoinHostPort(result.pgAddr, result.pgPort),
		result.PgBase)

	return &result
}
