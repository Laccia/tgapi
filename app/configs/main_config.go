package configs

import (
	"encoding/json"
	"fmt"

	"github.com/joho/godotenv"
)

type Server struct {
	PGconn *PGconn `json:"config_pg"`
}

func New() *Server {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	result := Server{

		PGconn: NewPG(),
	}

	tmp, err := json.MarshalIndent(result, "\t", " ")
	if err != nil {
		fmt.Println("Config", err)
	}

	fmt.Print(string(tmp) + "\n")

	return &result
}
