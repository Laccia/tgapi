package main

import (
	"tgapiV2/internal/app"
	"tgapiV2/internal/config"
	"tgapiV2/pkg"

	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.GetAppConfig()
	if err != nil || cfg == nil {
		panic(err)
	}

	log.Logger = pkg.NewLogger(cfg.LogCfg)

	log.Info().Str("comp:", "main").Msg("log initiated")
	//Start Application
	err = app.Run(cfg, log.Logger)

	if err != nil {
		log.Fatal().Err(err).Str("comp:", "main").Msg("can't run application, shutting down")
	}

}
